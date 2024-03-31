package servers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/theTardigrade/golang-basicServer/internal/middleware"
	"github.com/theTardigrade/golang-basicServer/internal/options"
)

const (
	httpDatumIndex = iota
	httpsDatumIndex
	totalData
)

type datum struct {
	server         *http.Server
	restartPending chan bool
	terminateDone  chan bool
	protocol       string
	port           int
	mutex          sync.RWMutex
}

var (
	data [totalData]*datum
	o    *options.Datum
)

func Init(opts *options.Datum) {
	o = opts

	for i := 0; i < totalData; i++ {
		data[i] = &datum{
			restartPending: make(chan bool, 1),
			terminateDone:  make(chan bool, 1),
		}
	}

	httpDatum, httpsDatum := data[httpDatumIndex], data[httpsDatumIndex]

	httpDatum.protocol = "http"
	httpsDatum.protocol = "https"

	httpDatum.port = opts.HttpPort
	httpsDatum.port = opts.HttpsPort

	httpDatum.server = &http.Server{
		// redirect all HTTP requests to corresponding route of HTTPS server
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hostname := strings.Split(r.Host, ":")[0]
			url := fmt.Sprintf("https://%s:%d%s", hostname, httpsDatum.port, r.URL.Path)
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		}),
		ReadHeaderTimeout: opts.RequestTimeoutDuration,
		WriteTimeout:      opts.ResponseTimeoutDuration,
	}

	httpsDatum.server = &http.Server{
		Handler:           middleware.Handler(opts),
		ReadHeaderTimeout: opts.RequestTimeoutDuration,
		WriteTimeout:      opts.ResponseTimeoutDuration,
	}
}
