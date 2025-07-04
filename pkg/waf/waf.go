package waf

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
	"net/netip"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/bloom42/stdx-go/httpx"
	"github.com/bloom42/stdx-go/log/slogx"
	"github.com/bloom42/stdx-go/memorycache"
	"github.com/bloom42/stdx-go/opt"
	"github.com/bloom42/stdx-go/retry"
	"github.com/bloom42/stdx-go/set"
	"github.com/tetratelabs/wazero"
	wazeroapi "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/experimental"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"markdown.ninja/pingoo-go"
	"markdown.ninja/pingoo-go/wasm"
	"markdown.ninja/pkg/server/httpctx"
)

type wasmModuleCtxKeyType struct{}

// wasmModuleCtxKeyType is the key that holds the wasmModule in a host WASM function call.
var wasmModuleCtxKey = wasmModuleCtxKeyType{}

type Waf struct {
	blockedCountries set.Set[string]
	logger           *slog.Logger
	dnsResolver      *net.Resolver

	wasmRuntime        wazero.Runtime
	compiledWasmModule wazero.CompiledModule
	wasmModulePool     *sync.Pool

	allowedBotIps *memorycache.Cache[netip.Addr, bool]
}

type wasmModule struct {
	module         wazeroapi.Module
	analyzeRequest wazeroapi.Function
	verifyBot      wazeroapi.Function
	allocate       wazeroapi.Function
	deallocate     wazeroapi.Function
}

var dnsServers = []string{
	"8.8.8.8:53",
	"1.0.0.1:53",
	"8.8.4.4:53",
	"1.1.1.1:53",
	// "9.9.9.9:53",
}

type analyzeRequestInput struct {
	HttpMethod       string     `json:"http_method"`
	UserAgent        string     `json:"user_agent"`
	IpAddress        netip.Addr `json:"ip_address"`
	Asn              int64      `json:"asn"`
	Country          string     `json:"country"`
	Path             string     `json:"path"`
	HttpVersionMajor int64      `json:"http_version_major"`
	HttpVersionMinor int64      `json:"http_version_minor"`
}

type outcome string

const (
	outcomeAllowed     outcome = "allowed"
	outcomeBlocked     outcome = "blocked"
	outcomeVerifiedBot outcome = "verified_bot"
)

type analyzeRequestOutput struct {
	Outcome outcome `json:"outcome"`
}

type lookupHostInput struct {
	IpAddress netip.Addr `json:"ip_address"`
	UserAgent string     `json:"user_agent"`
}

type lookupHostOutput struct {
	Hostname string `json:"hostname"`
}

type empty struct{}

func New(ctx context.Context, pingooClient *pingoo.Client, blockedCountries set.Set[string], logger *slog.Logger) (waf *Waf, err error) {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	allowedBotIps := memorycache.New(
		memorycache.WithTTL[netip.Addr, bool](7*24*time.Hour), // 7 days
		memorycache.WithCapacity[netip.Addr, bool](20_000),
	)

	logger.Debug("waf: downloading pingoo.wasm")
	pingooWasmBytes := bytes.NewBuffer(make([]byte, 0, 2_000_000)) // 2MB
	err = retry.Do(func() error {
		pingooWasmDownload, retryErr := pingooClient.DownloadPingooWasm(ctx, "")
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
		return nil, fmt.Errorf("waf: error downloading pingoo.wasm: %w", err)
	}

	logger.Debug("waf: pingoo.wasm successfully downloaded", slog.Int("size", pingooWasmBytes.Len()))

	wasmCtx := context.Background()
	// wasmCtx = experimental.WithMemoryAllocator(wasmCtx, wazeroallocator.NewNonMoving())

	// See https://github.com/tetratelabs/wazero/issues/2156
	// and https://github.com/wasilibs/go-re2/blob/main/internal/re2_wazero.go
	// for imformation about how to configure wazero to use a WASM lib using WASM memory

	// More wazero docs:
	// How to use HostFunctionBuilder with multiple goroutines? https://github.com/tetratelabs/wazero/issues/2217
	// Clarification on concurrency semantics for invocations https://github.com/tetratelabs/wazero/issues/2292
	// Improve InstantiateModule concurrency performance https://github.com/tetratelabs/wazero/issues/602
	// Add option to change Memory capacity https://github.com/tetratelabs/wazero/issues/500
	// Document best practices around invoking a wasi module multiple times https://github.com/tetratelabs/wazero/issues/985
	// API shape https://github.com/tetratelabs/wazero/issues/425

	wasmRuntime := wazero.NewRuntimeWithConfig(wasmCtx, wazero.NewRuntimeConfigCompiler().WithCoreFeatures(wazeroapi.CoreFeaturesV2|experimental.CoreFeaturesThreads).WithMemoryLimitPages(65536))

	wasi_snapshot_preview1.MustInstantiate(wasmCtx, wasmRuntime)

	// _, err = wasmRuntime.InstantiateWithConfig(wasmCtx, assets.MemoryWasm, wazero.NewModuleConfig().WithName("env"))
	// if err != nil {
	// 	return nil, fmt.Errorf("waf: error instantiating wasm memory module: %w", err)
	// }

	_, err = wasmRuntime.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithFunc(func(ctx context.Context, input wasm.Buffer) wasm.Buffer {
		return waf.wasmFunctionResolveHostForIp(ctx, input)
	}).Export("dns_lookup_ip_address").
		Instantiate(wasmCtx)
	if err != nil {
		return nil, fmt.Errorf("waf: error instantiating wasm host module (env): %w", err)
	}

	compiledWasmModule, err := wasmRuntime.CompileModule(wasmCtx, pingooWasmBytes.Bytes())
	if err != nil {
		return nil, fmt.Errorf("waf: error compiling wasm pingoo module: %w", err)
	}

	// as recommended in https://github.com/tetratelabs/wazero/issues/2217
	// we use a sync.Pool of wasm modules in order to handle concurrency
	wasmPool := &sync.Pool{
		New: func() any {
			poolObjectCtx := context.Background()
			instantiatedWasmModule, err := wasmRuntime.InstantiateModule(poolObjectCtx, compiledWasmModule, wazero.NewModuleConfig().
				WithStartFunctions("_initialize").WithSysNanosleep().WithSysNanotime().WithSysWalltime().WithName("").WithRandSource(cryptorand.Reader).WithStdout(os.Stdout).WithStderr(os.Stderr),
			// for debugging
			// .WithStdout(os.Stdout).WithStderr(os.Stderr),
			)
			if err != nil {
				logger.Error("waf.wasmModulePool.New: error instantiating WASM module", slogx.Err(err))
				return nil
			}

			poolObject := &wasmModule{
				module:         instantiatedWasmModule,
				analyzeRequest: instantiatedWasmModule.ExportedFunction("analyze_request"),
				verifyBot:      instantiatedWasmModule.ExportedFunction("verify_bot"),
				allocate:       instantiatedWasmModule.ExportedFunction("allocate"),
				deallocate:     instantiatedWasmModule.ExportedFunction("deallocate"),
			}
			// use a finalizer to Close the module, as recommended in https://github.com/golang/go/issues/23216
			runtime.SetFinalizer(poolObject, func(object *wasmModule) {
				object.module.Close(poolObjectCtx)
			})
			return poolObject
		},
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

	waf = &Waf{
		blockedCountries: blockedCountries,
		logger:           logger,
		allowedBotIps:    allowedBotIps,
		dnsResolver:      dnsResolver,

		wasmRuntime:        wasmRuntime,
		compiledWasmModule: compiledWasmModule,
		wasmModulePool:     wasmPool,
	}

	return
}

func (waf *Waf) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		userAgent := strings.TrimSpace(req.UserAgent())
		var err error
		path := req.URL.Path

		ctx := req.Context()
		httpCtx := httpctx.FromCtx(ctx)

		if len(userAgent) == 0 || len(userAgent) > 300 || !utf8.ValidString(userAgent) ||
			len(path) > 1024 || !utf8.ValidString(path) ||
			len(req.Method) > 20 ||
			waf.blockedCountries.Contains(httpCtx.Client.CountryCode) {
			waf.serveBlockedResponse(w)
			return
		}

		wasmModulePoolObject := waf.wasmModulePool.Get()
		if wasmModulePoolObject == nil {
			// fail open
			waf.logger.Error("waf: error getting object from wasm sync.Pool. Object is nil")
			next.ServeHTTP(w, req)
			return
		}

		wasmModule := wasmModulePoolObject.(*wasmModule)
		defer waf.wasmModulePool.Put(wasmModule)

		analyzeRequestAllowed, err := waf.analyzeRequest(req, wasmModule, userAgent)
		if err != nil {
			// fail open
			waf.logger.Error(err.Error(), slog.String("user_agent", userAgent),
				slog.String("ip_address", httpCtx.Client.IP.String()), slog.Int64("asn", httpCtx.Client.ASN))
			next.ServeHTTP(w, req)
			return
		}
		if !analyzeRequestAllowed {
			waf.serveBlockedResponse(w)
			return
		}

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}

// returns true if the request is allowed or false otherwise
func (waf *Waf) analyzeRequest(req *http.Request, wasmModule *wasmModule, userAgent string) (bool, error) {
	ctx := req.Context()
	httpCtx := httpctx.FromCtx(ctx)
	logger := slogx.FromCtx(ctx)

	ctx = context.WithValue(ctx, wasmModuleCtxKey, wasmModule)
	ctx = slogx.ToCtx(ctx, logger)

	analyzeRequestInputData := analyzeRequestInput{
		HttpMethod:       req.Method,
		UserAgent:        userAgent,
		IpAddress:        httpCtx.Client.IP,
		Asn:              httpCtx.Client.ASN,
		Country:          httpCtx.Client.CountryCode,
		Path:             req.URL.Path,
		HttpVersionMajor: int64(req.ProtoMajor),
		HttpVersionMinor: int64(req.ProtoMinor),
	}
	analyzeRequestRes, err := callWasmFunction[analyzeRequestInput, analyzeRequestOutput](ctx, wasmModule, wasmModule.analyzeRequest, analyzeRequestInputData)
	if err != nil {
		return false, fmt.Errorf("waf.analyzeRequest: error calling analyze_request wasm function: %w", err)
	}

	switch analyzeRequestRes.Outcome {
	case outcomeAllowed:
		return true, nil
	case outcomeVerifiedBot:
		waf.allowedBotIps.Set(httpCtx.Client.IP, true, memorycache.DefaultTTL)
		return true, nil
	case outcomeBlocked:
		return false, nil
	// case outcomeBot:
	default:
		// fail open
		waf.logger.Error("waf.analyzeRequest: unknown outcome", slog.String("outcome", string(analyzeRequestRes.Outcome)))
		return true, nil
	}
}

func (waf *Waf) serveBlockedResponse(res http.ResponseWriter) {
	sleepForMs := rand.Int64N(500) + 1000
	time.Sleep(time.Duration(sleepForMs) * time.Millisecond)

	message := "Access denied\n"

	res.Header().Set(httpx.HeaderConnection, "close")
	res.Header().Del(httpx.HeaderETag)
	res.Header().Set(httpx.HeaderCacheControl, httpx.CacheControlNoCache)
	res.Header().Set(httpx.HeaderContentType, httpx.MediaTypeTextUtf8)
	res.Header().Set(httpx.HeaderContentLength, strconv.FormatInt(int64(len(message)), 10))
	res.WriteHeader(http.StatusForbidden)
	res.Write([]byte(message))
}

func (waf *Waf) resolveHostForIp(ctx context.Context, ip netip.Addr) (string, error) {
	var hosts []string
	err := retry.Do(func() (retryErr error) {
		hosts, retryErr = waf.dnsResolver.LookupAddr(ctx, ip.String())
		if retryErr != nil {
			return retryErr
		}

		return nil
	}, retry.Context(ctx), retry.Attempts(4), retry.Delay(50*time.Millisecond))
	if err != nil {
		return "", fmt.Errorf("waf: error resolving hosts for IP address (%s): %w", ip, err)
	}

	cleanedUpHosts := make([]string, 0, len(hosts))
	for _, host := range hosts {
		host = strings.ToValidUTF8(strings.TrimSuffix(strings.TrimSpace(host), "."), "")
		if host != "" {
			cleanedUpHosts = append(cleanedUpHosts, host)
		}
	}
	hosts = cleanedUpHosts

	if len(hosts) > 0 {
		return hosts[0], nil
	}

	return "", nil
}

// callWasmFunction calls the given WASM function using JSON to serialize/deserialize input/output
func callWasmFunction[I, O any](ctx context.Context, wasmModule *wasmModule, wasmFunction wazeroapi.Function, input I) (O, error) {
	var emptyOutput O
	logger := slogx.FromCtx(ctx)

	// first we serialize the input ot JSON
	// then we allocate WASM memory for this JSON using the module's exported alloc function.
	// Don't forget to free the WASM input buffer
	// then we copy the input JSON into the WASM memoy
	// then we call the WASM function and pass it a pointer to the buffer that we have allocated
	// the function returns a pointer to a buffer it has allocated containing the output
	// then we read the output buffer from WASM's memory to the host (Go) memory and free the WASM ouput buffer
	// then we deserialize the output buffer content from JSON

	// serialize input to JSON
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return emptyOutput, fmt.Errorf("error marshalling input data to JSON: %w", err)
	}

	// allocate WASM memory for input data
	allocateInputResults, err := wasmModule.allocate.Call(ctx, uint64(len(inputBytes)))
	if err != nil {
		return emptyOutput, fmt.Errorf("error allocating wasm memory for function call input: %w", err)
	}

	wasmInputBuffer := wasm.Buffer(allocateInputResults[0])
	defer func() {
		// this memory was allocated by the WASM module so we have to deallocate it when finished
		if _, deallocateErr := wasmModule.deallocate.Call(ctx, uint64(wasmInputBuffer)); deallocateErr != nil {
			logger.Error("error deallocating wasm memory for function call input", slogx.Err(deallocateErr))
		}
	}()

	// write serialized input data into WASM's memory
	if !wasmModule.module.Memory().Write(wasmInputBuffer.Pointer(), inputBytes) {
		return emptyOutput, fmt.Errorf("error writing function call input data to wasm memory: Memory.Write(%d, %d) out of range of memory size %d",
			wasmInputBuffer.Pointer(), wasmInputBuffer.Size(), wasmModule.module.Memory().Size())
	}

	// call WASM function
	wasmResults, err := wasmFunction.Call(ctx, uint64(wasmInputBuffer))
	if err != nil {
		return emptyOutput, fmt.Errorf("error calling wasm function: %w", err)
	}

	// data returned from WASM's side is returned as an allocated buffer which (pointer, length) pair that is packed
	// into an uint64. It needs to be freed after having been read.
	wasmOutputBuffer := wasm.Buffer(wasmResults[0])
	defer func() {
		if _, deallocateErr := wasmModule.deallocate.Call(ctx, uint64(wasmOutputBuffer)); deallocateErr != nil {
			logger.Error("error deallocating wasm memory for function call output", slogx.Err(deallocateErr))
		}
	}()

	// read serialized output data from WASM memory
	outputBytes, outputReadOk := wasmModule.module.Memory().Read(wasmOutputBuffer.Pointer(), wasmOutputBuffer.Size())
	if !outputReadOk {
		return emptyOutput, fmt.Errorf("error reading function call output data from wasm memory: Memory.Read(%d, %d) out of range of memory size %d",
			wasmOutputBuffer.Pointer(), wasmOutputBuffer.Size(), wasmModule.module.Memory().Size())
	}

	var wasmResult wasm.Result[O]
	err = json.Unmarshal(outputBytes, &wasmResult)
	if err != nil {
		return emptyOutput, fmt.Errorf("error unmarshalling JSON output: %w", err)
	}

	if wasmResult.Error != nil {
		return emptyOutput, errors.New(*wasmResult.Error)
	}

	return *wasmResult.Ok, nil
}

func (waf *Waf) wasmFunctionResolveHostForIp(ctx context.Context, inputBuffer wasm.Buffer) wasm.Buffer {
	wasmModule := ctx.Value(wasmModuleCtxKey).(*wasmModule)

	inputBytes, readInputIp := wasmModule.module.Memory().Read(inputBuffer.Pointer(), inputBuffer.Size())
	if !readInputIp {
		return newWasmError(ctx, wasmModule, fmt.Errorf("error reading host function call input data from wasm memory: Memory.Read(%d, %d) out of range of memory size %d",
			inputBuffer.Pointer(), inputBuffer.Size(), wasmModule.module.Memory().Size()))
	}

	var input lookupHostInput
	err := json.Unmarshal(inputBytes, &input)
	if err != nil {
		return newWasmError(ctx, wasmModule, fmt.Errorf("error unmarshalling host function call input data from JSON: %w", err))
	}

	host, err := waf.resolveHostForIp(ctx, input.IpAddress)
	if err != nil {
		waf.logger.Warn(err.Error(), slog.String("user_agent", input.UserAgent))
		return newWasmError(ctx, wasmModule, err)
	}

	output := lookupHostOutput{
		Hostname: host,
	}
	outputResult := wasm.Result[lookupHostOutput]{
		Ok: &output,
	}
	outputBytes, err := json.Marshal(outputResult)
	if err != nil {
		return newWasmError(ctx, wasmModule, fmt.Errorf("error marshalling host function call output data to JSON: %w", err))
	}

	allocateOutputRes, err := wasmModule.allocate.Call(ctx, uint64(len(outputBytes)))
	if err != nil {
		return newWasmError(ctx, wasmModule, fmt.Errorf("error allocating memory for host function call output data: %w", err))
	}
	wasmOutputBuffer := wasm.Buffer(allocateOutputRes[0])

	if !wasmModule.module.Memory().Write(wasmOutputBuffer.Pointer(), outputBytes) {
		// TODO: log error?

		// deallocate := wazeroModule.ExportedFunction("deallocate")
		if _, deallocateErr := wasmModule.deallocate.Call(ctx, uint64(wasmOutputBuffer)); deallocateErr != nil {
			waf.logger.Error("error deallocating wasm memory for host function call output", slogx.Err(deallocateErr))
		}

		return newWasmError(ctx, wasmModule, fmt.Errorf("error writing host function call output data to wasm memory: Memory.Write(%d, %d) out of range of memory size %d",
			wasmOutputBuffer.Pointer(), wasmOutputBuffer.Size(), wasmModule.module.Memory().Size()))
	}

	return wasmOutputBuffer
}

func newWasmError(ctx context.Context, wasmModule *wasmModule, err error) wasm.Buffer {
	wasmErr := wasm.Result[empty]{
		Error: opt.String(err.Error()),
	}
	outputBytes, err := json.Marshal(wasmErr)
	if err != nil {
		// TODO: log error?
		return wasm.Buffer(0)
	}

	allocateOutputRes, err := wasmModule.allocate.Call(ctx, uint64(len(outputBytes)))
	if err != nil {
		// TODO: log error?
		return wasm.Buffer(0)
	}
	wasmOutputBuffer := wasm.Buffer(allocateOutputRes[0])

	if !wasmModule.module.Memory().Write(wasmOutputBuffer.Pointer(), outputBytes) {
		// TODO: log error?

		// deallocate := wazeroModule.ExportedFunction("deallocate")
		if _, deallocateErr := wasmModule.deallocate.Call(ctx, uint64(wasmOutputBuffer)); deallocateErr != nil {
			// TODO: log error?
			// waf.logger.Error("error deallocating wasm memory for host function call output", slogx.Err(deallocateErr))
		}

		return wasm.Buffer(0)
	}

	return wasmOutputBuffer
}
