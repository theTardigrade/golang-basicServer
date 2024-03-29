package options

import (
	"net/http"
	"time"
)

type Datum struct {
	HttpPort                int
	HttpsPort               int
	ResponseTimeoutDuration time.Duration
	RequestTimeoutDuration  time.Duration
	ShutdownTimeoutDuration time.Duration
	CertFilePath            string
	KeyFilePath             string
	Routes                  DatumRoutes
}

type DatumRoutes struct {
	Head   map[string]http.Handler
	Delete map[string]http.Handler
	Get    map[string]http.Handler
	Patch  map[string]http.Handler
	Put    map[string]http.Handler
	Post   map[string]http.Handler
}
