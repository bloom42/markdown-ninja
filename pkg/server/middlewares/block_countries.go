package middlewares

import (
	"net/http"
	"slices"

	"markdown.ninja/pkg/server/httpctx"
)

func BlockCountries(blockedCountries []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			httpCtx := httpctx.FromCtx(req.Context())
			if slices.Contains(blockedCountries, httpCtx.Client.CountryCode) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, req)
		}

		return http.HandlerFunc(fn)
	}
}
