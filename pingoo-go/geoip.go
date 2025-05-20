package pingoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/bloom42/stdx-go/httpx"
)

type GeoipDatabaseOutput struct {
	Data        io.ReadCloser
	NotModified bool
	Etag        string
}

func (client *Client) DownloadLatestGeoipDatabase(ctx context.Context, currentGeoipHashHex string) (ret GeoipDatabaseOutput, err error) {
	url := client.apiBaseUrl + "/downloads/geoip.mmdb"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "ApiKey "+client.apiKey)
	if currentGeoipHashHex != "" {
		req.Header.Add(httpx.HeaderIfNoneMatch, strconv.Quote(currentGeoipHashHex))
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("client.request: Doing HTTP request: %w", err)
		return
	}

	if res.StatusCode >= 400 {
		var body []byte

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			err = fmt.Errorf("reading body: %w", err)
			return
		}

		var apiErr ApiError
		err = json.Unmarshal(body, &apiErr)
		if err != nil {
			err = fmt.Errorf("decoding API response: %w", err)
			return
		}
		err = apiErr
		return
	}

	etag := res.Header.Get(httpx.HeaderETag)
	if etag != "" {
		unquotedEtag, unquoteErr := strconv.Unquote(etag)
		if unquoteErr == nil {
			etag = unquotedEtag
		}
	}

	ret = GeoipDatabaseOutput{
		Data:        res.Body,
		Etag:        etag,
		NotModified: false,
	}

	if res.StatusCode == 304 {
		_, _ = io.Copy(io.Discard, res.Body)
		ret.NotModified = true
	}

	return
}
