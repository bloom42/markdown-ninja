package pingoo

import (
	"bytes"
	"context"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"github.com/bloom42/stdx-go/countries"
	"github.com/bloom42/stdx-go/mmdb"
	"github.com/bloom42/stdx-go/retry"
)

const (
	AsUnknown string = "AS0"
)

type geoipDB struct {
	mmdbDatabase *mmdb.Reader
	// SHA-256 hash
	hash []byte
}

type GeoipRecord struct {
	// AsDomain string `maxminddb:"as_domain"`
	AsName string `json:"as_name"`
	ASNStr string `maxminddb:"asn"`
	ASN    int64  `maxminddb:"-"`
	// Continent     string `maxminddb:"continent"`
	// ContinentName string `maxminddb:"continent_name"`
	Country string `json:"country"`
	// CountryName   string `maxminddb:"country_name"`
}

// type geoipLookupInput struct {
// 	Ip string `json:"ip"`
// }

func (client *Client) GeoipLookup(ctx context.Context, ip netip.Addr) (ret GeoipRecord, err error) {
	// TODO: move this method to WASM.
	// but there is a problem: In order to perform geoip lookups, we need to load the geoip database in
	// (WASM) memory, making the WASM module stateful and use a "lot" ~50 MB of memory.
	// Thus, it's no longer possible to use pooling (via `sync.Pool`) to instantiate WASM module
	// for concurrent calls to guest functions (as recommended in https://github.com/tetratelabs/wazero/issues/2217
	// and https://github.com/tetratelabs/wazero/issues/985). Thus we need to compile the WASM module for the
	// `wasi-p1-threads` target to create a WASM module that supports concurrency and mutexes.
	// But then the module imports `memory` which is badly supported by wazero.
	// We got it half working following the instructions here: https://github.com/tetratelabs/wazero/issues/2156
	// but then we encountered errors such as `import memory[env.memory]: minimum size mismatch 23 > 17`
	// and `memory allocation of XX bytes failed`.

	// input := geoipLookupInput{
	// 	Ip: ip.String(),
	// }
	// ret, err = callWasmGuestFunction[geoipLookupInput, GeoipRecord](ctx, client, "geoip_lookup", input)
	// if err != nil {
	// 	ret.ASN = 0
	// 	ret.Country = countries.CodeUnknown
	// 	ret.AsName = ""
	// 	// don't return an error if IP not found
	// 	// TODO: better error detection
	// 	if !strings.Contains(err.Error(), "not found") {
	// 		return ret, err
	// 	}
	// 	err = nil
	// }

	err = client.geoipDB.Load().mmdbDatabase.Lookup(net.IP(ip.AsSlice()), &ret)
	if err != nil {
		ret.ASN = 0
		ret.ASNStr = AsUnknown
		ret.Country = countries.CodeUnknown
		err = fmt.Errorf("pingoo: error lookip up IP address [%s]: %w", ip, err)
		return
	}

	if ret.ASNStr == "" {
		ret.ASNStr = AsUnknown
	}
	asnInt, err := strconv.ParseInt(strings.TrimPrefix(ret.ASNStr, "AS"), 10, 64)
	if err != nil {
		err = fmt.Errorf("pingoo: error parsing ASN [%s]: %w", ret.ASNStr, err)
		return
	}
	if asnInt < 0 {
		err = fmt.Errorf("pingoo: error parsing ASN [%s]: ASN is negative", ret.ASNStr)
		return
	}
	ret.ASN = asnInt

	if ret.Country == "" {
		ret.Country = countries.CodeUnknown
	}

	return ret, nil
}

// type geoipSetDatabaseInput struct {
// 	Database []byte `json:"database"`
// }

// load the geoip database in the WASM module
// func (client *Client) geoipSetDatabase(ctx context.Context, database []byte) (ret wasm.Empty, err error) {
// 	input := geoipSetDatabaseInput{
// 		Database: database,
// 	}
// 	return callWasmGuestFunction[geoipSetDatabaseInput, wasm.Empty](ctx, client, "geoip_set_database", input)
// }

func (client *Client) refreshGeoipDatabase(ctx context.Context) (err error) {
	client.logger.Debug("pingoo: starting geoip database refresh")

	currentDatabase := client.geoipDB.Load()

	var currentDatabaseHashHex string
	if currentDatabase != nil {
		currentDatabaseHashHex = hex.EncodeToString(currentDatabase.hash)
	}

	geoipDbRes, err := client.DownloadLatestGeoipDatabase(ctx, currentDatabaseHashHex)
	if err != nil {
		return fmt.Errorf("pingoo: error downloading geoip database: %w", err)
	}
	defer geoipDbRes.Data.Close()

	if geoipDbRes.NotModified || geoipDbRes.Etag == currentDatabaseHashHex {
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

	geoipDBHash := sha3.Sum256(geoipDbBuffer.Bytes())
	// geoipDbHashHex := hex.EncodeToString(geoipDbRawHash)
	// if res.Etag != "" && res.Etag != geoipDbHashHex {
	// 	logger.Error("geoip: downloaded geoip database hash doesn't match etag",
	// 		slog.String("algorithm", "SHA3-256"), slog.String("encoding", "hex"),
	// 		slog.String("etag", res.Etag), slog.String("geoip_db.hash", geoipDbHashHex),
	// 	)
	// }

	mmdbDBReader, err := openAndValidateGeoipDatabase(geoipDbBuffer.Bytes())
	if err != nil {
		return err
	}

	geoipDb := &geoipDB{
		mmdbDatabase: mmdbDBReader,
		hash:         geoipDBHash[:],
	}
	client.geoipDB.Store(geoipDb)
	// _, err = client.geoipSetDatabase(ctx, geoipDbBuffer.Bytes())
	// if err != nil {
	// 	return fmt.Errorf("pingoo: error downloading geoip database: %w", err)
	// }

	client.logger.Info("geoip: geoip database successfully refreshed")

	return nil
}

type expectedGeoipResult struct {
	ip  string
	asn string
}

// make sure that the geoip database is valid
func openAndValidateGeoipDatabase(mmdbData []byte) (mmdbReader *mmdb.Reader, err error) {
	tests := []expectedGeoipResult{
		{ip: "1.1.1.1", asn: "13335"}, // Cloudflare
		{ip: "8.8.8.8", asn: "15169"}, // Google
	}

	// if database's size is < 10MB something is wrong
	if len(mmdbData) < 10_000_000 {
		err = fmt.Errorf("pingoo: geoip database is too small (%d)", len(mmdbData))
		return
	}

	// if database's size is > 200MB something is wrong
	if len(mmdbData) > 200_000_000 {
		err = fmt.Errorf("pingoo: geoip database is too big (%d)", len(mmdbData))
		return
	}

	mmdbReader, err = mmdb.FromBytes(mmdbData)
	if err != nil {
		err = fmt.Errorf("pingoo: error parsing mmdb file: %w", err)
		return
	}

	for _, test := range tests {
		var ip netip.Addr
		var ipInfo GeoipRecord

		ip, err = netip.ParseAddr(test.ip)
		if err != nil {
			err = fmt.Errorf("pingoo.openAndValidateGeoipDatabase: error parsing IP address [%s]: %w", test.ip, err)
			return
		}
		err = mmdbReader.Lookup(ip.AsSlice(), &ipInfo)
		if err != nil {
			err = fmt.Errorf("pingoo.openAndValidateGeoipDatabase: error looking up IP address [%s]: %w", test.ip, err)
			return
		}

		asn := strings.TrimPrefix(ipInfo.ASNStr, "AS")
		if asn != test.asn {
			err = fmt.Errorf("pingoo.openAndValidateGeoipDatabase: inconsistent ASN for IP address (%s): got: %s, expected: %s", test.ip, asn, test.asn)
			return
		}
	}

	return
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
