package pingoo

import (
	"context"
	"errors"

	"markdown.ninja/pingoo-go/wasm"
)

// callWasmGuestFunction calls the given WASM function using JSON to serialize/deserialize input/output
func callWasmGuestFunction[I, O any](ctx context.Context, client *Client, functionName string, parameters I) (O, error) {
	var emptyOutput O

	wasmModulePoolObject := client.wasmModulePool.Get()
	if wasmModulePoolObject == nil {
		err := errors.New("pingoo: error getting object from wasm sync.Pool. Object is nil")
		return emptyOutput, err
	}

	wasmModule := wasmModulePoolObject.(*wasm.Module)
	defer client.wasmModulePool.Put(wasmModule)

	output, err := wasm.CallGuestFunction[I, O](ctx, wasmModule, functionName, parameters)
	return output, err
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
