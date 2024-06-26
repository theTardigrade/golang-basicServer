package servers

import (
	"errors"
	"net"
	"time"
)

const (
	waitForOpenPortsTimeoutDuration = time.Second * 8
	waitForOpenPortsSleepDuration   = time.Millisecond * 100
)

var (
	ErrTimeout = errors.New("timeout waiting for server ports to open")
)

// WaitForOpenPorts makes sure that the web server ports are not currently in use
func WaitForOpenPorts() error {
	startTime := time.Now()

	for ; ; time.Sleep(waitForOpenPortsSleepDuration) {
		portsAvailable := true

		for _, datum := range data {
			var addr string

			func() {
				defer datum.mutex.RUnlock()
				datum.mutex.RLock()

				addr = datum.server.Addr
			}()

			ln, err := net.Listen("tcp", addr)

			if err != nil {
				portsAvailable = false
				break
			}

			ln.Close()
		}

		if portsAvailable {
			return nil
		}

		if time.Since(startTime) >= waitForOpenPortsTimeoutDuration {
			return ErrTimeout
		}
	}
}
