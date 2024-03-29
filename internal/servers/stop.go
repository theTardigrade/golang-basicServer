package servers

import (
	"context"
	"sync"

	"github.com/theTardigrade/golang-basicServer/internal/events"
)

func init() {
	events.AddNormalHandler(events.StopEvent, func() {
		var wg sync.WaitGroup

		for i := range data {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				stop(i, false)
			}(i)
		}

		wg.Wait()
	})
}

func stop(datumIndex int, restart bool) error {
	datum := data[datumIndex]

	datum.mutex.Lock()
	defer datum.mutex.Unlock()

	datum.restartPending <- restart

	// timeoutDuration, err := environment.Data.GetDuration("server_shutdown_timeout_duration")
	// if err != nil {
	// 	if environment.IsKeyNotFoundErr(err) {
	// 		timeoutDuration = time.Minute
	// 	} else {
	// 		panic(err)
	// 	}
	// }

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
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
