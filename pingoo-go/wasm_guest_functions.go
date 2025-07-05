package pingoo

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

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

type verifyJwtInput struct {
	JWKS  Jwks   `json:"jwks"`
	Token string `json:"token"`
}

type verifyJwtOutput[C any] struct {
	Claims C `json:"claims"`
	// header
	// signature
}

func VerifyJWT[C any](ctx context.Context, client *Client, token string) (claims C, err error) {
	wasmInput := verifyJwtInput{
		JWKS:  client.jwks,
		Token: token,
	}
	wasmOutput, err := callWasmGuestFunction[verifyJwtInput, verifyJwtOutput[C]](ctx, client, "verify_jwt", wasmInput)
	if err != nil {
		return claims, fmt.Errorf("pingoo.VerifyJWT: error calling verify_jwt wasm function: %w", err)
	}

	return wasmOutput.Claims, nil
}

type AnalyzeRequestInput struct {
	HttpMethod       string     `json:"http_method"`
	UserAgent        string     `json:"user_agent"`
	IpAddress        netip.Addr `json:"ip_address"`
	Asn              int64      `json:"asn"`
	Country          string     `json:"country"`
	Path             string     `json:"path"`
	HttpVersionMajor int64      `json:"http_version_major"`
	HttpVersionMinor int64      `json:"http_version_minor"`
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

func AnalyzeRequest(ctx context.Context, client *Client, input AnalyzeRequestInput) (ret AnalyzeRequestOutput, err error) {
	return callWasmGuestFunction[AnalyzeRequestInput, AnalyzeRequestOutput](ctx, client, "analyze_request", input)
}
