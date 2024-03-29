package routes

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/theTardigrade/golang-basicServer/internal/router"
)

var (
	staticGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var success bool

		defer func() {
			if err := recover(); err != nil {
				panic(err)
			} else if !success {
				router.Multiplexer.HandleNotFound(w, r)
			}
		}()

		localPath := filepath.Join(
			staticGlobalFilepathBase,
			prevDirRegexp.ReplaceAllLiteralString(r.URL.Path[1:], prevDirReplacement),
		)

		fileInfo, err := os.Stat(localPath)
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}

			return
		}
		if !fileInfo.Mode().IsRegular() {
			return
		}
		fileModTime := fileInfo.ModTime()

		fileContents, err := ioutil.ReadFile(localPath)
		if err != nil {
			panic(err)
		}

		header := w.Header()

		header.Set("Content-Type", datum.mimeType)

		useEtag := true

		if fileInfo.Size() < staticEtagFileSizeMin {
			useEtag = false
		}

		success = true

		if useEtag {
			etag := datum.etag

			if r.Header.Get("If-None-Match") == etag {
				if gw, ok := w.(*grw.GzipResponseWriter); ok {
					gw.UnsetHeaders()
				}

				w.WriteHeader(http.StatusNotModified)
				w.Write([]byte{})
				return
			}

			header.Set("Etag", etag)
		}

		w.WriteHeader(http.StatusOK)
		w.Write(datum.fileContents)
	})
)

const (
	staticPath = "/static/*"
)

func init() {
	router.Multiplexer.GetFunc(staticPath, staticGetHandler)
}
