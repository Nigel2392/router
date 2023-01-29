package middleware

import (
	"compress/gzip"
	"net/http"

	"github.com/Nigel2392/router"
)

func GZIP(next router.Handler) router.Handler {
	return router.HandleFuncWrapper{F: func(v router.Vars, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		// Compress the response
		var gz = gzip.NewWriter(w)
		defer gz.Close()
		// Create gzip response writer
		var gzw = gzipResponseWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(v, gzw, r)
	}}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) { return w.Writer.Write(b) }
func (w gzipResponseWriter) Header() http.Header         { return w.ResponseWriter.Header() }
