package pingoo

import (
	"context"
	"net/http"
	"strings"
)

type analyzeRequestInput struct {
	HttpMethod       string       `json:"http_method"`
	Hostname         string       `json:"hostname"`
	UserAgent        string       `json:"user_agent"`
	Ip               string       `json:"ip"`
	Asn              int64        `json:"asn"`
	Country          string       `json:"country"`
	Path             string       `json:"path"`
	HttpVersionMajor int64        `json:"http_version_major"`
	HttpVersionMinor int64        `json:"http_version_minor"`
	Headers          []httpHeader `json:"headers"`
}

// Serializing HTTP headers as []structure is significantly more efficient than using slices
// BenchmarkSerializeHttpHeader/struct-5         	  614881	      2088 ns/op	    1496 B/op	       3 allocs/op
// BenchmarkSerializeHttpHeader/slice-5          	  372471	      3252 ns/op	    2025 B/op	      45 allocs/op
type httpHeader struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type AnalyzeRequestOutcome string

const (
	AnalyzeRequestOutcomeAllowed     AnalyzeRequestOutcome = "allowed"
	AnalyzeRequestOutcomeBlocked     AnalyzeRequestOutcome = "blocked"
	AnalyzeRequestOutcomeVerifiedBot AnalyzeRequestOutcome = "verified_bot"
	AnalyzeRequestOutcomeChallenge   AnalyzeRequestOutcome = "challenge"
)

type analyzeRequestOutput struct {
	Outcome AnalyzeRequestOutcome `json:"outcome"`
}

func convertHttpheaders(headers http.Header) []httpHeader {
	ret := make([]httpHeader, 0, len(headers))
	for headerName, headerValues := range headers {
		ret = append(ret, httpHeader{Name: strings.ToLower(headerName), Values: headerValues})
	}
	return ret
}

func (client *Client) analyzeRequest(ctx context.Context, input analyzeRequestInput) (ret analyzeRequestOutput, err error) {
	analyzeRequestRes, err := callWasmGuestFunction[analyzeRequestInput, analyzeRequestOutput](ctx, client, "analyze_request", input)
	if err != nil {
		return
	}

	switch analyzeRequestRes.Outcome {
	case AnalyzeRequestOutcomeVerifiedBot:
		return client.verifyBot(ctx, input)
	default:
		return analyzeRequestRes, nil
	}
}

type verifyBotInput struct {
	HttpMethod string `json:"http_method"`
	UserAgent  string `json:"user_agent"`
	Ip         string `json:"ip"`
	Asn        int64  `json:"asn"`
	Path       string `json:"path"`
	IpHostname string `json:"ip_hostname"`
}

func (client *Client) verifyBot(ctx context.Context, input analyzeRequestInput) (ret analyzeRequestOutput, err error) {
	ipHostnameRes, err := client.resolveHostForIp(ctx, lookupHostInput{IpAddress: input.Ip})
	if err != nil {
		return
	}

	verifyBotInputData := verifyBotInput{
		HttpMethod: input.HttpMethod,
		UserAgent:  input.UserAgent,
		Ip:         input.Ip,
		Asn:        input.Asn,
		Path:       input.Path,
		IpHostname: ipHostnameRes.Hostname,
	}
	return callWasmGuestFunction[verifyBotInput, analyzeRequestOutput](ctx, client, "verify_bot", verifyBotInputData)
}
