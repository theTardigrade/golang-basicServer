package middleware

import (
	"net/http"

	"github.com/theTardigrade/golang-basicServer/internal/options"
	"github.com/theTardigrade/golang-basicServer/internal/router"
)

func Handler(opts *options.Datum) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l := len(opts.Middleware.Before); l > 0 {
			h := opts.Middleware.Before[l-1](router.Handler)

			for i := l - 2; i >= 0; i-- {
				h = opts.Middleware.Before[i](h)
			}

			h.ServeHTTP(w, r)
		} else {
			router.Handler.ServeHTTP(w, r)
		}

		if l := len(opts.Middleware.After); l > 0 {
			var h http.Handler = emptyHandler

			for i := l - 1; i >= 0; i-- {
				h = opts.Middleware.After[i](h)
			}

			h.ServeHTTP(w, r)
		}
	})
}
