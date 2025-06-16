package geoip

import (
	"context"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bloom42/stdx-go/log/slogx"
	"github.com/bloom42/stdx-go/mmdb"
	"github.com/bloom42/stdx-go/retry"
	"markdown.ninja/pingoo-go"
)

const (
	CountryCodeUnknown string = "XX"
	AsUnknown          string = "AS0"
)

type Record struct {
	// AsDomain string `maxminddb:"as_domain"`
	AsName string `maxminddb:"as_name"`
	ASNStr string `maxminddb:"asn"`
	ASN    int64  `maxminddb:"-"`
	// Continent     string `maxminddb:"continent"`
	// ContinentName string `maxminddb:"continent_name"`
	Country string `maxminddb:"country"`
	// CountryName   string `maxminddb:"country_name"`
}

type Resolver struct {
	pingooClient *pingoo.Client
	logger       *slog.Logger
	geoipDB      atomic.Pointer[database]
}

type database struct {
	mmdbDatabase *mmdb.Reader
	// SHA-512 hash
	hash []byte
}

func Init(ctx context.Context, pingooClient *pingoo.Client, logger *slog.Logger) (resolver *Resolver, err error) {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	resolver = &Resolver{
		pingooClient: pingooClient,
		logger:       logger,
		geoipDB:      atomic.Pointer[database]{},
	}

	err = retry.Do(func() error {
		retryErr := resolver.DownloadLatestGeoipDatabase(ctx)
		if retryErr != nil {
			resolver.logger.Warn("geoip: error downloading geoip database", slogx.Err(retryErr))
		}
		return retryErr
	}, retry.Context(ctx), retry.Attempts(15), retry.Delay(time.Second), retry.DelayType(retry.FixedDelay))
	if err != nil {
		return nil, err
	}

	return resolver, nil
}

func (resolver *Resolver) Lookup(ip netip.Addr) (info Record, err error) {
	err = resolver.geoipDB.Load().mmdbDatabase.Lookup(net.IP(ip.AsSlice()), &info)
	if err != nil {
		info.ASN = 0
		info.ASNStr = AsUnknown
		info.Country = CountryCodeUnknown
		err = fmt.Errorf("geoip: error lookip up IP address [%s]: %w", ip, err)
		return
	}

	if info.ASNStr == "" {
		info.ASNStr = AsUnknown
	}
	asnInt, err := strconv.ParseInt(strings.TrimPrefix(info.ASNStr, "AS"), 10, 64)
	if err != nil {
		err = fmt.Errorf("geoip: error parsing ASN [%s]: %w", info.ASNStr, err)
		return
	}
	if asnInt < 0 {
		err = fmt.Errorf("geoip: error parsing ASN [%s]: ASN is negative", info.ASNStr)
		return
	}
	info.ASN = asnInt

	if info.Country == "" {
		info.Country = CountryCodeUnknown
	}

	return
}

func IsPrivate(ip netip.Addr) bool {
	return ip.IsInterfaceLocalMulticast() || ip.IsLinkLocalMulticast() ||
		ip.IsLoopback() || ip.IsMulticast() || ip.IsPrivate() || ip.IsUnspecified()
}

func (resolver *Resolver) DownloadLatestGeoipDatabase(ctx context.Context) (err error) {
	logger := slogx.FromCtx(ctx)

	currentDatabase := resolver.geoipDB.Load()
	var currentDatabaseHashHex string
	if currentDatabase != nil {
		currentDatabaseHash := resolver.geoipDB.Load().hash
		currentDatabaseHashHex = hex.EncodeToString(currentDatabaseHash)
	}

	res, err := resolver.pingooClient.DownloadLatestGeoipDatabase(ctx, currentDatabaseHashHex)
	if err != nil {
		return fmt.Errorf("geoip: fetching latest geoip database: %w", err)
	}
	defer res.Data.Close()

	if res.NotModified || res.Etag == currentDatabaseHashHex {
		resolver.logger.Info("geoip: no new geoip database is available")
		return nil
	}

	tmpFile, err := os.CreateTemp("", "markdown_ninja_tmp_geoip_database")
	if err != nil {
		return fmt.Errorf("geoip: error creating tmp file: %w", err)
	}
	// best effort cleanup
	defer os.Remove(tmpFile.Name())

	geoipDbHasher := sha3.New256()
	dataHasherReader := io.TeeReader(res.Data, geoipDbHasher)
	_, err = io.Copy(tmpFile, dataHasherReader)
	_ = tmpFile.Close()
	if err != nil && err != io.EOF {
		return fmt.Errorf("geoip: writing geoip database to tmp file: %w", err)
	}
	err = nil

	geoipDbRawHash := geoipDbHasher.Sum(nil)
	geoipDbHashHex := hex.EncodeToString(geoipDbRawHash)
	if res.Etag != "" && res.Etag != geoipDbHashHex {
		logger.Error("geoip: downloaded geoip database hash doesn't match etag",
			slog.String("algorithm", "SHA3-256"), slog.String("encoding", "hex"),
			slog.String("etag", res.Etag), slog.String("geoip_db.hash", geoipDbHashHex),
		)
	}

	mmdbData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("geoip: reading downloaded geoip database: %w", err)
	}

	mmdbDBReader, err := openAndValidateGeoipDatabase(mmdbData)
	if err != nil {
		return err
	}

	geoipDb := &database{
		mmdbDatabase: mmdbDBReader,
		hash:         geoipDbRawHash,
	}
	resolver.geoipDB.Store(geoipDb)

	resolver.logger.Info("geoip: new geoip database successfully downloaded")

	return nil
}

type expectedGeoipResult struct {
	ip  string
	asn string
}

func openAndValidateGeoipDatabase(mmdbData []byte) (mmdbReader *mmdb.Reader, err error) {
	tests := []expectedGeoipResult{
		{ip: "1.1.1.1", asn: "13335"}, // Cloudflare
		{ip: "8.8.8.8", asn: "15169"}, // Google
	}

	// if database's size is < 10MB something is wrong
	if len(mmdbData) < 10_000_000 {
		err = fmt.Errorf("geoip: database is too small (%d)", len(mmdbData))
		return
	}

	// if database's size is > 200MB something is wrong
	if len(mmdbData) > 200_000_000 {
		err = fmt.Errorf("geoip: database is too big (%d)", len(mmdbData))
		return
	}

	mmdbReader, err = mmdb.FromBytes(mmdbData)
	if err != nil {
		err = fmt.Errorf("geoip: parsing mmdb file: %w", err)
		return
	}

	for _, test := range tests {
		var ip netip.Addr
		var ipInfo Record

		ip, err = netip.ParseAddr(test.ip)
		if err != nil {
			err = fmt.Errorf("geoip.openAndValidateGeoipDatabase: error parsing IP address [%s]: %w", test.ip, err)
			return
		}
		err = mmdbReader.Lookup(ip.AsSlice(), &ipInfo)
		if err != nil {
			err = fmt.Errorf("geoip.openAndValidateGeoipDatabase: error looking up IP address [%s]: %w", test.ip, err)
			return
		}

		asn := strings.TrimPrefix(ipInfo.ASNStr, "AS")
		if asn != test.asn {
			err = fmt.Errorf("geoip: inconsistent ASN for IP address (%s): got: %s, expected: %s", test.ip, asn, test.asn)
			return
		}
	}

	return
}
