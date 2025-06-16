package middleware

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"markdown.ninja/pingoo-go"
)

// LoggingMiddleware is a middleware that logs HTTP request information.
func LoggingMiddleware(ctx context.Context, pingooClient *pingoo.Client, config pingoo.MiddlewareConfig) func(next http.Handler) http.Handler {
	logsBuffer := newlogsBuffer(pingooClient)
	go logsBuffer.flushInBackground(ctx)
	defaultLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return func(next http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			if config.Logging.Disabled {
				next.ServeHTTP(res, req)
				return
			}

			var logger *slog.Logger

			// Create a response writer that captures the status code and response size
			lrw := &LoggingResponseWriter{ResponseWriter: res}

			// Call the next handler
			startTime := time.Now()
			next.ServeHTTP(lrw, req)
			durationMs := time.Since(startTime).Milliseconds()

			_, clientIP, _ := extractClientIpAddress(req)

			hostname := req.Host
			responseSize := lrw.size
			if responseSize > math.MaxInt64 {
				responseSize = math.MaxInt64
			}

			userAgent := req.UserAgent()
			if !utf8.ValidString(userAgent) {
				userAgent = strings.ToValidUTF8(userAgent, "")
			}

			path := req.URL.Path
			if !utf8.ValidString(path) {
				path = strings.ToValidUTF8(path, "")
			}

			// Log the request information
			logRecord := pingoo.HttpLogRecord{
				Time:         startTime,
				Duration:     durationMs,
				Method:       req.Method,
				Path:         path,
				Host:         hostname,
				ClientIP:     clientIP,
				UserAgent:    userAgent,
				StatusCode:   lrw.statusCode,
				ResponseSize: int64(responseSize),
				HTTPVersion:  req.Proto,
			}
			logsBuffer.Push(logRecord)

			if config.Logging.GetLogger != nil {
				logger = config.Logging.GetLogger(req.Context())
			} else {
				logger = defaultLogger
			}

			if logger != nil {
				logger.Info("HTTP request",
					slog.String("method", logRecord.Method),
					slog.String("path", logRecord.Path),
					slog.String("host", logRecord.Host),
					slog.Int64("duration", durationMs),
					slog.String("client_ip", logRecord.ClientIP.String()),
					slog.String("user_agent", logRecord.UserAgent),
					slog.Uint64("status_code", uint64(logRecord.StatusCode)),
					slog.Uint64("response_size", uint64(logRecord.ResponseSize)),
					slog.String("http_version", logRecord.HTTPVersion),
				)
			}
		}

		return http.HandlerFunc(fn)
	}
}

// LoggingResponseWriter is a custom ResponseWriter that captures the status code and response size.
// a better version is available for inspiration at https://github.com/go-chi/chi/blob/master/middleware/wrap_writer.go
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode uint16
	size       uint64
	hijacked   bool
}

// WriteHeader captures the status code.
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	if code >= 0 && code <= math.MaxUint16 {
		lrw.statusCode = uint16(code)
	} else {
		lrw.statusCode = math.MaxUint16
	}

	lrw.ResponseWriter.WriteHeader(code)
}

// Write captures the response size.
func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.hijacked {
		return 0, http.ErrHijacked
	}
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += uint64(size)
	return size, err
}

func (lrw *LoggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	conn, rw, err := hijacker.Hijack()
	if err == nil {
		lrw.hijacked = true
	}
	return conn, rw, err
}

func (lrw *LoggingResponseWriter) Flush() {
	flusher, ok := lrw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

// http.CloseNotifier is deprecated:
// func (lrw *LoggingResponseWriter) CloseNotify() <-chan bool {
// 	return lrw.ResponseWriter.(http.CloseNotifier).CloseNotify()
// }

// func (lrw *LoggingResponseWriter) Pusher() (pusher http.Pusher) {
// 	if pusher, ok := lrw.ResponseWriter.(http.Pusher); ok {
// 		return pusher
// 	}
// 	return nil
// }

func extractClientIpAddress(req *http.Request) (clientIpStr string, clientIp netip.Addr, err error) {
	clientIpStr, _, err = net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		err = fmt.Errorf("extractClientIpAddress: RemoteAddr (%s) is not valid: %w", clientIpStr, err)
		return
	}

	clientIp, err = netip.ParseAddr(clientIpStr)
	if err != nil {
		err = fmt.Errorf("extractClientIpAddress: error parsing client IP [%s]: %w", clientIpStr, err)
		return
	}

	return
}
