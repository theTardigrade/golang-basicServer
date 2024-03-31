package middleware

import "net/http"

var (
	emptyHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)
