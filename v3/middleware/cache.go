package middleware

import (
	"strconv"

	"github.com/Nigel2392/router/v3"
	"github.com/Nigel2392/router/v3/request"
)

// Set the cache headers for the response.
// This will enable caching for the specified amount of seconds.
func Cache(maxAge int) func(next router.Handler) router.Handler {
	return func(next router.Handler) router.Handler {
		return router.HandleFunc(func(r *request.Request) {
			for _, header := range etagHeaders {
				r.Response.Header().Del(header)
			}
			r.Response.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
			next.ServeHTTP(r)
		})
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

// Set the cache headers for the response.
// This will disable caching.
func NoCache(next router.Handler) router.Handler {
	return router.HandleFunc(func(r *request.Request) {
		for _, header := range etagHeaders {
			r.Response.Header().Del(header)
		}
		r.Response.Header().Set("Cache-Control", "no-cache, no-store, no-transform, must-revalidate, private, max-age=0")
		r.Response.Header().Set("Pragma", "no-cache")
		r.Response.Header().Set("Expires", "0")
		next.ServeHTTP(r)
	})
}
