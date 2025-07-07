package pingoo

import (
	"context"
	"net/http"
	"time"

	"github.com/bloom42/stdx-go/log/slogx"
	"markdown.ninja/pingoo-go/wasm"
)

// callWasmGuestFunction calls the given WASM function using JSON to serialize/deserialize input/output
func callWasmGuestFunction[I, O any](ctx context.Context, client *Client, functionName string, parameters I) (O, error) {
	logger := slogx.FromCtx(ctx)
	logger.Debug("wasm: calling function: " + functionName)
	functionStartedAt := time.Now()

	// TODO: make this call able to handle concurrency without mutex.
	// For that, we first need to compile the WASM module to target `wasm32-wasip1-threds`.
	// Then we need to take in account wazero concurrency support (https://github.com/tetratelabs/wazero/issues/2217)
	// especially related to exported function calls (https://pkg.go.dev/github.com/tetratelabs/wazero/api#Function).
	// Unfortunately, as of now, we couldn't make it work. WASM calls don't work as if something in the stack
	// was not handling concurrency correctly, but we were unable to track down what.
	//
	// error calling wasm function: wasm error: out of bounds memory access
	//
	// The other option is to use a `sync.Pool` of wasm module, but then the modules need to be stateless
	// otherwise they will use too much memory (1 geoip database in memory per WASM module) and we need to
	// push more logic in the Go SDK instead of the WASM SDK.

	// var emptyOutput O
	// wasmModulePoolObject := client.wasmModulePool.Get()
	// if wasmModulePoolObject == nil {
	// 	err := errors.New("pingoo: error getting object from wasm sync.Pool. Object is nil")
	// 	return emptyOutput, err
	// }

	// wasmModule := wasmModulePoolObject.(*wasm.Module)
	// defer client.wasmModulePool.Put(wasmModule)
	wasmModule := client.wasmModule
	client.wasmModuleMutex.Lock()
	defer client.wasmModuleMutex.Unlock()

	output, err := wasm.CallGuestFunction[I, O](ctx, wasmModule, functionName, parameters)
	if err != nil {
		logger.Debug("wasm: function " + functionName + " completed in " + time.Since(functionStartedAt).String() + " with error: " + err.Error())
		return output, err
	}

	logger.Debug("wasm: function " + functionName + " successfully completed in " + time.Since(functionStartedAt).String())

	return output, nil
}

type AnalyzeRequestInput struct {
	HttpMethod       string `json:"http_method"`
	UserAgent        string `json:"user_agent"`
	IpAddress        string `json:"ip_address"`
	Asn              int64  `json:"asn"`
	Country          string `json:"country"`
	Path             string `json:"path"`
	HttpVersionMajor int64  `json:"http_version_major"`
	HttpVersionMinor int64  `json:"http_version_minor"`
}

type AnalyzeRequestOutcome string

const (
	AnalyzeRequestOutcomeAllowed     AnalyzeRequestOutcome = "allowed"
	AnalyzeRequestOutcomeBlocked     AnalyzeRequestOutcome = "blocked"
	AnalyzeRequestOutcomeVerifiedBot AnalyzeRequestOutcome = "verified_bot"
)

type AnalyzeRequestOutput struct {
	Outcome AnalyzeRequestOutcome `json:"outcome"`
}

func (client *Client) AnalyzeRequest(ctx context.Context, input AnalyzeRequestInput) (ret AnalyzeRequestOutput, err error) {
	return callWasmGuestFunction[AnalyzeRequestInput, AnalyzeRequestOutput](ctx, client, "analyze_request", input)
}

type handleHttpRequestInput struct {
	Path string `json:"path"`
}

type handleHttpRequestOutput struct {
	Status uint16 `json:"status"`
	// Vec<[2]string>
	Headers [][]string `json:"headers"`
	Body    []byte     `json:"body"`
}

func (client *Client) handleHttpRequest(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	input := handleHttpRequestInput{
		Path: req.URL.Path,
	}
	result, err := callWasmGuestFunction[handleHttpRequestInput, handleHttpRequestOutput](ctx, client, "handle_http_request", input)
	if err != nil {
		client.serveBlockedResponse(res)
		return
	}

	for _, header := range result.Headers {
		res.Header().Set(header[0], header[1])
	}

	res.WriteHeader(int(result.Status))
	res.Write(result.Body)
}
