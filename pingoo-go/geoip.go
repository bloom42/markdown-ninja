package pingoo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/netip"
	"strings"
	"time"

	"github.com/bloom42/stdx-go/countries"
	"github.com/bloom42/stdx-go/opt"
	"github.com/bloom42/stdx-go/retry"
	"markdown.ninja/pingoo-go/wasm"
)

const (
	AsUnknown string = "AS0"
)

type GeoipRecord struct {
	// AsDomain string `maxminddb:"as_domain"`
	AsName string `json:"as_name"`
	ASN    int64  `json:"-"`
	// Continent     string `maxminddb:"continent"`
	// ContinentName string `maxminddb:"continent_name"`
	Country string `json:"country"`
	// CountryName   string `maxminddb:"country_name"`
}

type geoipLookupInput struct {
	Ip string `json:"ip"`
}

func (client *Client) GeoipLookup(ctx context.Context, ip netip.Addr) (ret GeoipRecord, err error) {
	input := geoipLookupInput{
		Ip: ip.String(),
	}
	ret, err = callWasmGuestFunction[geoipLookupInput, GeoipRecord](ctx, client, "geoip_lookup", input)
	if err != nil {
		ret.ASN = 0
		ret.Country = countries.CodeUnknown
		ret.AsName = ""
		// don't return an error if IP not found
		// TODO: better error detection
		if !strings.Contains(err.Error(), "not found") {
			return ret, err
		}
		err = nil
	}

	return ret, nil
}

type geoipSetDatabaseInput struct {
	Database []byte `json:"database"`
}

// loads the geoip database in the WASM module
func (client *Client) geoipSetDatabase(ctx context.Context, database []byte) (ret wasm.Empty, err error) {
	input := geoipSetDatabaseInput{
		Database: database,
	}
	return callWasmGuestFunction[geoipSetDatabaseInput, wasm.Empty](ctx, client, "geoip_set_database", input)
}

func (client *Client) refreshGeoipDatabase(ctx context.Context) (err error) {
	client.logger.Debug("pingoo: starting geoip database refresh")

	currentGeoipDBEtag := client.geoipDBEtag.Load()
	if currentGeoipDBEtag == nil {
		currentGeoipDBEtag = opt.String("")
	}

	geoipDbRes, err := client.DownloadLatestGeoipDatabase(ctx, *currentGeoipDBEtag)
	if err != nil {
		return fmt.Errorf("pingoo: error downloading geoip database: %w", err)
	}
	defer geoipDbRes.Data.Close()

	if geoipDbRes.NotModified || geoipDbRes.Etag == *currentGeoipDBEtag {
		client.logger.Debug("pingoo: no new geoip database is available")
		return nil
	}

	// download the actual database in a buffer
	geoipDbBuffer := bytes.NewBuffer(make([]byte, 0, 50_000_000))
	_, err = io.Copy(geoipDbBuffer, geoipDbRes.Data)
	if err != nil && err != io.EOF {
		return fmt.Errorf("pingoo: error downloading geoip database: %w", err)
	}
	err = nil

	// if res.Etag != "" && res.Etag != geoipDbHashHex {
	// 	logger.Error("geoip: downloaded geoip database hash doesn't match etag",
	// 		slog.String("algorithm", "SHA3-256"), slog.String("encoding", "hex"),
	// 		slog.String("etag", res.Etag), slog.String("geoip_db.hash", geoipDbHashHex),
	// 	)
	// }

	_, err = client.geoipSetDatabase(ctx, geoipDbBuffer.Bytes())
	if err != nil {
		return fmt.Errorf("pingoo: error downloading geoip database: %w", err)
	}

	client.geoipDBEtag.Store(&geoipDbRes.Etag)

	client.logger.Info("geoip: geoip database successfully refreshed")

	return nil
}

func (client *Client) refreshGeoipDatabaseInBackground(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(client.geoipDBRefreshInterval):
		}

		err := retry.Do(func() error {
			retryErr := client.refreshGeoipDatabase(ctx)
			return retryErr
		}, retry.Context(ctx), retry.Attempts(15), retry.Delay(time.Minute), retry.DelayType(retry.FixedDelay))
		if err != nil {
			client.logger.Warn("pingoo: error refreshing geoip database: " + err.Error())
		}
	}
}

// func IsPrivate(ip netip.Addr) bool {
// 	return ip.IsInterfaceLocalMulticast() || ip.IsLinkLocalMulticast() ||
// 		ip.IsLoopback() || ip.IsMulticast() || ip.IsPrivate() || ip.IsUnspecified()
// }
