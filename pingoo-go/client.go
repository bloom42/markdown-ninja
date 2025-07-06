package pingoo

import (
	"bytes"
	"context"
	cryptorand "crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bloom42/stdx-go/httpx"
	"github.com/bloom42/stdx-go/log/slogx"
	"github.com/bloom42/stdx-go/retry"
	"github.com/bloom42/stdx-go/uuid"
	"github.com/klauspost/compress/zstd"
	"github.com/tetratelabs/wazero"
	wazeroapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"markdown.ninja/pingoo-go/assets"
	"markdown.ninja/pingoo-go/wasm"
	"markdown.ninja/pkg/jwt"
)

const userAgent = "Pingoo/GoSDK (https://pingoo.io)"

var dnsServers = []string{
	"8.8.8.8:53",
	"1.0.0.1:53",
	"8.8.4.4:53",
	"1.1.1.1:53",
	// "9.9.9.9:53",
}

type ClientConfig struct {
	Url        *string
	HttpClient *http.Client
	Logger     *slog.Logger
}

type Client struct {
	pingooURL  string
	projectId  uuid.UUID
	apiKey     string
	apiBaseUrl string
	httpClient *http.Client
	logger     *slog.Logger

	dnsResolver *net.Resolver

	jwksRefreshInterval time.Duration
	jwks                jwt.Jwks
	jwksLock            sync.RWMutex

	wasmRuntime        wazero.Runtime
	compiledWasmModule wazero.CompiledModule
	wasmModulePool     *sync.Pool

	geoipDBRefreshInterval time.Duration
	geoipDB                atomic.Pointer[geoipDB]
	// wasmModule *wasm.Module
	// wasmMutex  sync.Mutex
}

type ApiError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (err ApiError) Error() string {
	return err.Message
}

// TODO: wrap errors with errs
func NewClient(ctx context.Context, apiKey string, projectID string, config *ClientConfig) (client *Client, err error) {
	url := "https://pingoo.io"
	httpClient := httpx.DefaultClient()

	if config != nil {
		if config.Url != nil && *config.Url != "" {
			url = *config.Url
		}
		if config.HttpClient != nil {
			httpClient = config.HttpClient
		}
	}

	projectId, err := uuid.Parse(projectID)
	if err != nil {
		return nil, errors.New("pingoo: projectID is not valid")
	}

	logger := config.Logger
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	dnsResolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{
				Timeout: 5 * time.Second,
			}
			dnsServer := dnsServers[rand.IntN(len(dnsServers))]
			return dialer.DialContext(ctx, network, dnsServer)
		},
	}

	client = &Client{
		pingooURL:           url,
		apiBaseUrl:          url + "/api",
		projectId:           projectId,
		apiKey:              apiKey,
		httpClient:          httpClient,
		logger:              logger,
		jwksRefreshInterval: time.Minute,
		// jwksKeys:          make(map[string]VerifyingKey),
		// jwks:               Jwks{Keys: []Jwk{}},
		jwksLock:           sync.RWMutex{},
		dnsResolver:        dnsResolver,
		wasmRuntime:        nil,
		compiledWasmModule: nil,

		geoipDBRefreshInterval: 12 * time.Hour,
		geoipDB:                atomic.Pointer[geoipDB]{},
		// wasmModule:         nil,
		// wasmMutex:          sync.Mutex{},
	}

	// JWKS

	err = retry.Do(func() error {
		retryErr := client.refreshJwks(ctx)
		return retryErr
	}, retry.Context(ctx), retry.Attempts(15), retry.Delay(time.Second), retry.DelayType(retry.FixedDelay))
	if err != nil {
		return nil, fmt.Errorf("pingoo: error downloading JWKS: %w", err)
	}

	go client.refreshJwksInBackground(ctx)

	// WASM

	logger.Debug("pingoo: downloading pingoo.wasm")
	pingooWasmBytes := bytes.NewBuffer(make([]byte, 0, 2_000_000)) // 2MB
	err = retry.Do(func() error {
		pingooWasmDownload, retryErr := client.DownloadPingooWasm(ctx, "")
		if retryErr != nil {
			return fmt.Errorf("downloading file: %w", retryErr)
		}
		defer pingooWasmDownload.Data.Close()
		_, retryErr = io.Copy(pingooWasmBytes, pingooWasmDownload.Data)
		if retryErr != nil {
			return fmt.Errorf("copying bytes: %w", retryErr)
		}
		return nil
	}, retry.Context(ctx), retry.Attempts(15), retry.Delay(time.Second), retry.DelayType(retry.FixedDelay))
	if err != nil {
		return nil, fmt.Errorf("pingoo: error downloading pingoo.wasm: %w", err)
	}

	logger.Debug("pingoo: pingoo.wasm successfully downloaded", slog.Int("size", pingooWasmBytes.Len()))

	wasmCtx := context.Background()

	// See https://github.com/tetratelabs/wazero/issues/2156
	// and https://github.com/wasilibs/go-re2/blob/main/internal/re2_wazero.go
	// for imformation about how to configure wazero to use a WASM lib using WASM memory
	// wasmCtx = experimental.WithMemoryAllocator(wasmCtx, wazeroallocator.NewNonMoving())

	// More wazero docs:
	// How to use HostFunctionBuilder with multiple goroutines? https://github.com/tetratelabs/wazero/issues/2217
	// Clarification on concurrency semantics for invocations https://github.com/tetratelabs/wazero/issues/2292
	// Improve InstantiateModule concurrency performance https://github.com/tetratelabs/wazero/issues/602
	// Add option to change Memory capacity https://github.com/tetratelabs/wazero/issues/500
	// Document best practices around invoking a wasi module multiple times https://github.com/tetratelabs/wazero/issues/985
	// API shape https://github.com/tetratelabs/wazero/issues/425

	client.wasmRuntime = wazero.NewRuntimeWithConfig(wasmCtx, wazero.NewRuntimeConfigCompiler().WithCoreFeatures(wazeroapi.CoreFeaturesV2|experimental.CoreFeaturesThreads).WithMemoryLimitPages(65536))

	wasi_snapshot_preview1.MustInstantiate(wasmCtx, client.wasmRuntime)

	_, err = client.wasmRuntime.InstantiateWithConfig(wasmCtx, assets.MemoryWasm, wazero.NewModuleConfig().WithName("env"))
	if err != nil {
		return nil, fmt.Errorf("pingoo: error instantiating wasm memory module: %w", err)
	}

	_, err = client.wasmRuntime.NewHostModuleBuilder("host").
		NewFunctionBuilder().WithFunc(func(ctx context.Context, input wasm.Buffer) wasm.Buffer {
		return wasm.HandleHostFunctionCall(ctx, client.resolveHostForIp, input)
	}).Export("dns_lookup_ip_address").
		Instantiate(wasmCtx)
	if err != nil {
		return nil, fmt.Errorf("pingoo: error instantiating wasm host module (host): %w", err)
	}

	client.compiledWasmModule, err = client.wasmRuntime.CompileModule(wasmCtx, pingooWasmBytes.Bytes())
	if err != nil {
		return nil, fmt.Errorf("pingoo: error compiling wasm pingoo module: %w", err)
	}

	// instantiatedWasmModule, err := client.wasmRuntime.InstantiateModule(wasmCtx, client.compiledWasmModule, wazero.NewModuleConfig().
	// 	WithStartFunctions("_initialize").WithSysNanosleep().WithSysNanotime().WithSysWalltime().WithName("").WithRandSource(cryptorand.Reader).WithStdout(os.Stdout).WithStderr(os.Stderr),
	// // for debugging
	// // .WithStdout(os.Stdout).WithStderr(os.Stderr),
	// )
	// if err != nil {
	// 	return nil, fmt.Errorf("pingoo: error instantiating WASM module: %w", err)
	// }

	// client.wasmModule, err = wasm.NewModule(instantiatedWasmModule)
	// if err != nil {
	// 	return nil, fmt.Errorf("pingoo: error instantiating WASM module: %w", err)
	// }

	// runtime.SetFinalizer(client.wasmModule, func(module *wasm.Module) {
	// 	module.Close(wasmCtx)
	// })

	// as recommended in https://github.com/tetratelabs/wazero/issues/2217
	// we use a sync.Pool of wasm modules in order to handle concurrent calls of WASM functions
	client.wasmModulePool = &sync.Pool{
		New: func() any {
			poolObjectCtx := context.Background()
			instantiatedWasmModule, err := client.wasmRuntime.InstantiateModule(poolObjectCtx, client.compiledWasmModule, wazero.NewModuleConfig().
				WithStartFunctions("_initialize").WithSysNanosleep().WithSysNanotime().WithSysWalltime().WithName("").WithRandSource(cryptorand.Reader).WithStdout(os.Stdout).WithStderr(os.Stderr),
			// for debugging
			// .WithStdout(os.Stdout).WithStderr(os.Stderr),
			)
			if err != nil {
				logger.Error("waf.wasmModulePool.New: error instantiating WASM module", slogx.Err(err))
				return nil
			}

			poolObject, err := wasm.NewModule(instantiatedWasmModule)
			if err != nil {
				logger.Error("waf.wasmModulePool.New: error instantiating WASM module", slogx.Err(err))
				return nil
			}

			// use a finalizer to Close the module, as recommended in https://github.com/golang/go/issues/23216
			runtime.SetFinalizer(poolObject, func(module *wasm.Module) {
				module.Close(poolObjectCtx)
			})

			return poolObject
		},
	}

	// GeoIP

	err = retry.Do(func() error {
		retryErr := client.refreshGeoipDatabase(ctx)
		return retryErr
	}, retry.Context(ctx), retry.Attempts(15), retry.Delay(time.Second), retry.DelayType(retry.FixedDelay))
	if err != nil {
		return nil, fmt.Errorf("pingoo: error downloading geoip database: %w", err)
	}

	go client.refreshGeoipDatabaseInBackground(ctx)

	return client, nil
}

type requestParams struct {
	Method  string
	Route   string
	Payload any
}

func (client *Client) request(ctx context.Context, params requestParams, dst any) error {
	url := client.apiBaseUrl + params.Route

	req, err := http.NewRequestWithContext(ctx, params.Method, url, nil)
	if err != nil {
		return err
	}

	if params.Payload != nil {
		jsonPayload, err := json.Marshal(params.Payload)
		if err != nil {
			return fmt.Errorf("client.request: marshaling JSON: %w", err)
		}

		// zstdCompressor, err := zstd.NewWriter(nil,
		// 	zstd.WithEncoderCRC(true),
		// 	zstd.WithEncoderConcurrency(1),
		// )
		// if err != nil {
		// 	return fmt.Errorf("client.request: error instantiating zstd encoder: %w", err)
		// }

		// compressedPayload := zstdCompressor.EncodeAll(jsonPayload, nil)
		// zstdCompressor.Close()

		// req.Body = io.NopCloser(bytes.NewReader(compressedPayload))
		req.Body = io.NopCloser(bytes.NewReader(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Encoding", "zstd")
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "zstd")
	req.Header.Set("Authorization", "ApiKey "+client.apiKey)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("client.request: Doing HTTP request: %w", err)
	}
	defer res.Body.Close()

	var body []byte

	resContentEncoding := res.Header.Get("Content-Encoding")
	if resContentEncoding != "" {
		if resContentEncoding != "zstd" {
			return fmt.Errorf("client.request: Content-Encoding \"%s\" not supported. Only zstd is currently supported", resContentEncoding)
		}

		bodyBuffer := bytes.NewBuffer(make([]byte, 0, 1000))

		zstdDecompressor, err := zstd.NewReader(res.Body,
			zstd.WithDecoderConcurrency(1),
		)
		if err != nil {
			return fmt.Errorf("client.request: error instantiating zstd encoder: %w", err)
		}

		_, err = io.Copy(bodyBuffer, zstdDecompressor)
		if err != nil {
			zstdDecompressor.Close()
			return fmt.Errorf("client.request: error decompressing response's body: %w", err)
		}
		zstdDecompressor.Close()
		body = bodyBuffer.Bytes()
	} else {
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("client.request: Reading body: %w", err)
		}
	}

	if res.StatusCode >= 400 {
		var apiErr ApiError
		err = json.Unmarshal(body, &apiErr)
		if err != nil {
			return fmt.Errorf("pingoo: error decoding API response: %w", err)
		}
		return apiErr
	}

	if dst != nil {
		err = json.Unmarshal(body, &dst)
		if err != nil {
			return fmt.Errorf("pingoo: error decoding API response: %w", err)
		}
	}

	return nil
}
