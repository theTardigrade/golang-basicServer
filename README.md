# golang-basicServer

This package makes it easy to create an HTTPS server (alongside an HTTP server that will redirect to the more secure one).

## Example

```golang
package main

import (
	"net/http"
	"time"

	basicServer "github.com/theTardigrade/golang-basicServer"
)

const (
	output = "<!DOCTYPE><html><body><h1>THIS IS AN EXAMPLE</h1></body></html>"
)

var (
	opts = basicServer.Options{
		HttpPort:                80,
		HttpsPort:               443,
		ResponseTimeoutDuration: 60 * time.Second,
		RequestTimeoutDuration:  60 * time.Second,
		CertFilePath:            "/example/certificates/default.cer",
		KeyFilePath:             "/example/certificates/default.key",
		Routes: basicServer.OptionsRoutes{
			Get: map[string]http.Handler{
				"/": http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
					resp.Header().Set("Content-Type", "text/html")
					resp.WriteHeader(http.StatusOK)
					resp.Write([]byte(output))
				}),
			},
		},
	}
)

func main() {
	basicServer.ServeContinuously(opts, func(err error) {
		log.Println(err)
	})
}
```