package servers

import (
	"context"
)

func stop(datumIndex int, restart bool) error {
	datum := data[datumIndex]

	datum.mutex.Lock()
	defer datum.mutex.Unlock()

	datum.restartPending <- restart

	var ctx context.Context
	var cancel context.CancelFunc

	if o.ShutdownTimeoutDuration == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), waitForOpenPortsTimeoutDuration)
	}
	defer cancel()

	datum.server.Shutdown(ctx)

	if !restart {
		<-datum.terminateDone // wait for start function to return
	}

	return ctx.Err()
}

// StopHTTP shuts down the HTTP web server
func StopHTTP(restart bool) error {
	return stop(httpDatumIndex, restart)
}

// StopHTTPS shuts down the HTTPS web server
func StopHTTPS(restart bool) error {
	return stop(httpsDatumIndex, restart)
}
