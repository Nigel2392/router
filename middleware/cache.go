package middleware

import (
	"net/http"
	"strconv"

	"github.com/Nigel2392/router"
)

func Cache(maxAge int) func(next router.Handler) router.Handler {
	return func(next router.Handler) router.Handler {
		return router.HandleFuncWrapper{F: func(v router.Vars, w http.ResponseWriter, r *http.Request) {
			for _, header := range etagHeaders {
				w.Header().Del(header)
			}
			w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
			next.ServeHTTP(v, w, r)
		}}
	}
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(next router.Handler) router.Handler {
	return router.HandleFuncWrapper{F: func(v router.Vars, w http.ResponseWriter, r *http.Request) {
		for _, header := range etagHeaders {
			w.Header().Del(header)
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(v, w, r)
	}}
}
