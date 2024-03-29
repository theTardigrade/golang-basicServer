package router

import (
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/theTardigrade/golang-basicServer/internal/options"
)

var (
	Multiplexer = bone.New()
	Handler     = http.HandlerFunc(Multiplexer.ServeHTTP)
)

func Init(opts *options.Datum) {
	for path, handler := range opts.Routes.Head {
		Multiplexer.Head(path, handler)
	}

	for path, handler := range opts.Routes.Delete {
		Multiplexer.Delete(path, handler)
	}

	for path, handler := range opts.Routes.Get {
		Multiplexer.Get(path, handler)
	}

	for path, handler := range opts.Routes.Patch {
		Multiplexer.Patch(path, handler)
	}

	for path, handler := range opts.Routes.Put {
		Multiplexer.Put(path, handler)
	}

	for path, handler := range opts.Routes.Post {
		Multiplexer.Post(path, handler)
	}
}

func IsHandled(r *http.Request) bool {
	for _, routes := range Multiplexer.Routes {
		for _, route := range routes {
			if route.Match(r) {
				return true
			}
		}
	}

	return false
}
