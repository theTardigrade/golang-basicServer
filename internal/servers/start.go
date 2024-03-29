package servers

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

func start(datumIndex int, fatalErrChan chan<- error) {
	datum := data[datumIndex]

	runningChan := make(chan bool)
	errChan := make(chan error)

	datum.mutex.Lock()

	go func(ec chan<- error, rc chan<- bool) {
		var err error

		go func() {
			time.Sleep(time.Microsecond * 100)
			rc <- true
		}()

		switch datumIndex {
		case httpsDatumIndex:
			err = datum.server.ListenAndServeTLS(
				o.CertFilePath,
				o.KeyFilePath,
			)
		default:
			err = datum.server.ListenAndServe()
		}

		if err != nil {
			ec <- err
		}
	}(errChan, runningChan)

	port, protocol := datum.port, datum.protocol

	<-runningChan
	datum.mutex.Unlock()

	log.Printf("%s server started on port %d\n", strings.ToUpper(protocol), port)

	func(ec <-chan error, msg string) {
		if err := <-ec; err != nil {
			msg += fmt.Sprintf(" (%s)", err)
		}
		log.Println(msg)
	}(errChan, strings.ToUpper(protocol)+" server stopped")

	select {
	case restart := <-datum.restartPending:
		if restart {
			start(datumIndex, fatalErrChan)
			return
		}
	default:
	}

	datum.terminateDone <- true
	fatalErrChan <- errors.New(protocol + " server terminated")
}

// StartHTTP sets the HTTP web server running
func StartHTTP(fatalErrChan chan<- error) {
	start(httpDatumIndex, fatalErrChan)
}

// StartHTTPS sets the HTTPS web server running
func StartHTTPS(fatalErrChan chan<- error) {
	start(httpsDatumIndex, fatalErrChan)
}
