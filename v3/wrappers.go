//go:build !js && !wasm
// +build !js,!wasm

package router

import (
	"net/http"

	"github.com/Nigel2392/router/v3/request"
	"github.com/Nigel2392/router/v3/request/writer"
)

// Middleware is the function that is called when a route is matched
type Middleware func(Handler) Handler

type ErrorHandlerMiddleware func(error, *request.Request) Middleware

// Handler is the interface that wraps the ServeHTTP method.
type Handler interface {
	ServeHTTP(*request.Request)
}

// HandleFunc is the function that is called when a route is matched
type HandleFunc func(*request.Request)

func (f HandleFunc) ServeHTTP(r *request.Request) {
	f(r)
}

// Wrapper function for http.Handler to make it compatible with Handler
type httpHandlerWrapper struct {
	H http.Handler
}

// ServeHTTP implements the Handler interface
func (h httpHandlerWrapper) ServeHTTP(r *request.Request) {
	h.H.ServeHTTP(r.Response, r.Request)
}

// Make a new handler from a http.Handler
func FromHTTPHandler(h http.Handler) Handler {
	return httpHandlerWrapper{H: h}
}

// Make a new http.Handler from a Handler
func ToHTTPHandler(h Handler) http.Handler {
	var f = func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(request.NewRequest(writer.NewClearable(w), r, nil))
	}
	return http.HandlerFunc(f)
}
