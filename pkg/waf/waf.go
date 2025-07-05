package waf

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/bloom42/stdx-go/httpx"
	"github.com/bloom42/stdx-go/memorycache"
	"markdown.ninja/pingoo-go"
	"markdown.ninja/pkg/server/httpctx"
)

type Waf struct {
	logger *slog.Logger

	pingooClient *pingoo.Client

	allowedBotIps *memorycache.Cache[netip.Addr, bool]
}

func New(ctx context.Context, pingooClient *pingoo.Client, logger *slog.Logger) (waf *Waf, err error) {
	if logger == nil {
		logger = slog.New(slog.DiscardHandler)
	}

	allowedBotIps := memorycache.New(
		memorycache.WithTTL[netip.Addr, bool](7*24*time.Hour), // 7 days
		memorycache.WithCapacity[netip.Addr, bool](20_000),
	)

	waf = &Waf{
		logger:        logger,
		allowedBotIps: allowedBotIps,
		pingooClient:  pingooClient,
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
			len(req.Method) > 20 {
			waf.serveBlockedResponse(w)
			return
		}

		analyzeRequestInput := pingoo.AnalyzeRequestInput{
			HttpMethod:       req.Method,
			UserAgent:        userAgent,
			IpAddress:        httpCtx.Client.IP,
			Asn:              httpCtx.Client.ASN,
			Country:          httpCtx.Client.CountryCode,
			Path:             req.URL.Path,
			HttpVersionMajor: int64(req.ProtoMajor),
			HttpVersionMinor: int64(req.ProtoMinor),
		}
		analyzeRequestOutput, err := pingoo.AnalyzeRequest(ctx, waf.pingooClient, analyzeRequestInput)
		if err != nil {
			// fail open
			waf.logger.Error(err.Error(), slog.String("user_agent", userAgent),
				slog.String("ip_address", httpCtx.Client.IP.String()), slog.Int64("asn", httpCtx.Client.ASN))
			next.ServeHTTP(w, req)
			return
		}

		switch analyzeRequestOutput.Outcome {
		case pingoo.AnalyzeRequestOutcomeBlocked:
			waf.serveBlockedResponse(w)
			return
		case pingoo.AnalyzeRequestOutcomeAllowed:
			break
		case pingoo.AnalyzeRequestOutcomeVerifiedBot:
			waf.allowedBotIps.Set(httpCtx.Client.IP, true, memorycache.DefaultTTL)
		default:
			// fail open
			waf.logger.Error("waf.analyzeRequest: unknown outcome", slog.String("outcome", string(analyzeRequestOutput.Outcome)))
		}

		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
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
