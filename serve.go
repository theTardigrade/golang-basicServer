package basicServer

import (
	"github.com/theTardigrade/golang-basicServer/internal/options"
	"github.com/theTardigrade/golang-basicServer/internal/router"
	"github.com/theTardigrade/golang-basicServer/internal/servers"
)

type (
	Options       = options.Datum
	OptionsRoutes = options.DatumRoutes
)

func Serve(opts Options) (err error) {
	options.Init(&opts)
	router.Init(&opts)
	servers.Init(&opts)

	if err = servers.WaitForOpenPorts(); err != nil {
		return
	}

	fatalErrHTTPChan := make(chan error)
	fatalErrHTTPSChan := make(chan error)

	go servers.StartHTTPS(fatalErrHTTPSChan)
	go servers.StartHTTP(fatalErrHTTPChan)

	select {
	case err = <-fatalErrHTTPSChan:
		servers.StopHTTPS(true)
	case err = <-fatalErrHTTPChan:
		servers.StopHTTP(true)
	}

	return
}
