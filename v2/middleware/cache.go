package middleware

import (
	"strconv"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

func Cache(maxAge int) func(next router.Handler) router.Handler {
	return func(next router.Handler) router.Handler {
		return router.HandleFuncWrapper{F: func(r *request.Request) {
			for _, header := range etagHeaders {
				r.Writer.Header().Del(header)
			}
			r.Writer.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
			next.ServeHTTP(r)
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
	return router.HandleFuncWrapper{F: func(r *request.Request) {
		for _, header := range etagHeaders {
			r.Writer.Header().Del(header)
		}
		r.Writer.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
		r.Writer.Header().Set("Pragma", "no-cache")
		r.Writer.Header().Set("Expires", "0")
		next.ServeHTTP(r)
	}}
}
