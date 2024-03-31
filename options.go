package basicServer

import "github.com/theTardigrade/golang-basicServer/internal/options"

type (
	Options                  = options.Datum
	OptionsRoutes            = options.DatumRoutes
	OptionsMiddleware        = options.DatumMiddleware
	OptionsMiddlewareHandler = options.DatumMiddlewareHandler
)
