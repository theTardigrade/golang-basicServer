package options

import (
	"net/http"
	"time"
)

type Datum struct {
	HttpPort                int
	HttpsPort               int
	RequestTimeoutDuration  time.Duration
	ResponseTimeoutDuration time.Duration
	ShutdownTimeoutDuration time.Duration
	CertFilePath            string
	KeyFilePath             string
	Routes                  DatumRoutes
	Middleware              DatumMiddleware
}

type DatumRoutes struct {
	NotFound http.Handler
	Head     map[string]http.Handler
	Delete   map[string]http.Handler
	Get      map[string]http.Handler
	Patch    map[string]http.Handler
	Put      map[string]http.Handler
	Post     map[string]http.Handler
}

type DatumMiddleware struct {
	Before []DatumMiddlewareHandler
	After  []DatumMiddlewareHandler
}

type DatumMiddlewareHandler func(http.Handler) http.Handler

const (
	datumHttpPortDefault                = 80
	datumHttpsPortDefault               = 443
	datumRequestTimeoutDurationDefault  = 60 * time.Second
	datumResponseTimeoutDurationDefault = 60 * time.Second
	datumShutdownTimeoutDurationDefault = 60 * time.Second
)

const (
	datumPortMin     = 1
	datumPortMax     = 65535
	datumDurationMin = time.Nanosecond
)

func Init(opts *Datum) {
	if opts.HttpPort < datumPortMin || opts.HttpPort > datumPortMax {
		opts.HttpPort = datumHttpPortDefault
	}

	if opts.HttpsPort < datumPortMin || opts.HttpsPort > datumPortMax {
		opts.HttpsPort = datumHttpPortDefault
	}

	if opts.RequestTimeoutDuration < datumDurationMin {
		opts.RequestTimeoutDuration = datumRequestTimeoutDurationDefault
	}

	if opts.ResponseTimeoutDuration < datumDurationMin {
		opts.ResponseTimeoutDuration = datumResponseTimeoutDurationDefault
	}

	if opts.ShutdownTimeoutDuration < datumDurationMin {
		opts.ShutdownTimeoutDuration = datumShutdownTimeoutDurationDefault
	}
}
