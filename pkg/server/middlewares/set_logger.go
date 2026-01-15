package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/skerkour/stdx-go/log/slogx"
	"github.com/skerkour/stdx-go/uuid"
)

// SetLogger injects `logger` in the context of requests
func SetLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			requestLogger := logger

			reqIDContextValue := ctx.Value(RequestIDCtxKey)
			if requestID, ok := reqIDContextValue.(uuid.UUID); ok {
				requestLogger = requestLogger.With(slog.String("request_id", requestID.String()))
			}

			req = req.WithContext(slogx.ToCtx(ctx, requestLogger))

			next.ServeHTTP(res, req)
		}
		return http.HandlerFunc(fn)
	}
}
