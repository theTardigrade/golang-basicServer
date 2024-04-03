package basicServer

import (
	"net/http"

	"github.com/go-zoo/bone"
)

func ValueFromRequest(r *http.Request, key string) string {
	return bone.GetValue(r, key)
}

func ValuesMapFromRequest(r *http.Request) map[string]string {
	return bone.GetAllValues(r)
}
