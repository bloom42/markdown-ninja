package pingoo

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/netip"
)

type handleHttpRequestInput struct {
	Path string `json:"path"`
	Body []byte `json:"body"`
}

type handleHttpRequestOutput struct {
	Status uint16 `json:"status"`
	// Vec<[2]string>
	Headers [][]string `json:"headers"`
	Body    []byte     `json:"body"`
}

func (client *Client) handleHttpRequest(ctx context.Context, ip netip.Addr, res http.ResponseWriter, req *http.Request) {
	client.logger.Debug("pingoo.handleHttpRequest: " + req.URL.Path)

	switch req.URL.Path {
	case "/__pingoo/challenge/api/init":
		client.handleChallengeInitRequest(ctx, ip, res, req)
		return
	case "/__pingoo/challenge/api/verify":
		client.handleChallengeVerifyRequest(ctx, ip, res, req)
		return
	}

	bodyBuffer := bytes.NewBuffer(make([]byte, 0, 2_000))
	limitedBodyReader := http.MaxBytesReader(res, req.Body, 5_000_000) // 5 MB
	_, err := io.Copy(bodyBuffer, limitedBodyReader)
	if err != nil {
		// TODO: correct error handling
		client.serveInternalError(res)
		return
	}

	input := handleHttpRequestInput{
		Path: req.URL.Path,
		Body: bodyBuffer.Bytes(),
	}
	result, err := callWasmGuestFunction[handleHttpRequestInput, handleHttpRequestOutput](ctx, client, "handle_http_request", input)
	if err != nil {
		client.serveInternalError(res)
		return
	}

	for _, header := range result.Headers {
		res.Header().Add(header[0], header[1])
	}

	res.WriteHeader(int(result.Status))
	res.Write(result.Body)
}
