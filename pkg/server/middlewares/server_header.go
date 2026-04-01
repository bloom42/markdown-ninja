package middlewares

import (
	"net/http"

	"github.com/skerkour/stdx-go/httpx"
)

// SetServerHeader set's the Server HTTP header to "Markdown Ninja"
func SetServerHeader(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Server", "Markdown Ninja")
		w.Header().Add(httpx.HeaderAltSvc, `h3=":443"; ma=86400`)
		next.ServeHTTP(w, req)
	}

	return http.HandlerFunc(fn)
}
