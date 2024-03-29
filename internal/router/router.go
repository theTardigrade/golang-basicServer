package router

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/theTardigrade/golang-basicServer/internal/options"
)

var (
	multiplexer = bone.New()
	Handler     = http.HandlerFunc(multiplexer.ServeHTTP)
)

func Init(opts *options.Datum) {
	for path, handler := range opts.Routes.Get {
		multiplexer.Get(path, handler)
	}
}

func IsHandled(r *http.Request) bool {
	for _, routes := range multiplexer.Routes {
		for _, route := range routes {
			if route.Match(r) {
				return true
			}
		}
	}

	return false
}
