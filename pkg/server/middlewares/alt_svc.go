package middlewares

import (
	"net/http"

	"github.com/skerkour/stdx-go/httpx"
)

func AddAltSvcHeader() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(res, req)
			res.Header().Add(httpx.HeaderAltSvc, `h3=":443"; ma=86400`)
		}

		return http.HandlerFunc(fn)
	}
}
