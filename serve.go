package basicServer

import (
	"github.com/theTardigrade/golang-basicServer/internal/options"
	"github.com/theTardigrade/golang-basicServer/internal/router"
	"github.com/theTardigrade/golang-basicServer/internal/servers"
)

func Serve(opts Options) (err error) {
	serveInit(&opts)

	err = serve(&opts)

	return
}

func ServeContinuously(opts Options, errHandler func(error)) {
	serveInit(&opts)

	if errHandler == nil {
		for {
			serve(&opts)
		}
	} else {
		for {
			if err := serve(&opts); err != nil {
				errHandler(err)
			}
		}
	}
}

func serveInit(opts *Options) {
	options.Init(opts)
	router.Init(opts)
	servers.Init(opts)
}

func serve(opts *Options) (err error) {
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
