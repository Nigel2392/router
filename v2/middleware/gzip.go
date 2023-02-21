package middleware

import (
	"compress/gzip"
	"net/http"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

func GZIP(next router.Handler) router.Handler {
	return router.HandleFuncWrapper{F: func(r *request.Request) {
		r.Writer.Header().Set("Content-Encoding", "gzip")
		// Compress the response
		var gz = gzip.NewWriter(r.Writer)
		defer gz.Close()
		// Create gzip response writer
		var gzw = gzipResponseWriter{ResponseWriter: r.Writer, Writer: gz}
		r.Writer = gzw
		next.ServeHTTP(r)
	}}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) { return w.Writer.Write(b) }
func (w gzipResponseWriter) Header() http.Header         { return w.ResponseWriter.Header() }