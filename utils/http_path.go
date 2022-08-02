package utils

import "net/http"

func HttpPath(r *http.Request) string {
	if r.URL.RawPath != "" {
		return r.URL.RawPath
	} else {
		return r.URL.Path
	}
}
