package router

import (
	"net/http"

	"github.com/Nigel2392/router/v2/request"
)

// Make a new Handler
func ToHandler(f func(*request.Request)) Handler {
	return handleFuncWrapper{F: f}
}

// Make a new handler from a http.Handler
func FromHTTPHandler(h http.Handler) Handler {
	return httpHandlerWrapper{H: h}
}

// Make a new http.Handler from a Handler
func ToHTTPHandler(h Handler) http.Handler {
	var f = func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(request.NewRequest(w, r, nil))
	}
	return http.HandlerFunc(f)
}

// Wrapper function for http.Handler to make it compatible with Handler
type httpHandlerWrapper struct {
	H http.Handler
}

// ServeHTTP implements the Handler interface
func (h httpHandlerWrapper) ServeHTTP(r *request.Request) {
	h.H.ServeHTTP(r.Response, r.Request)
}

// Wrapper function for HandleFunc to make it compatible with http.Handler
type handleFuncWrapper struct {
	F func(*request.Request)
}

// ServeHTTP implements the Handler interface
func (h handleFuncWrapper) ServeHTTP(r *request.Request) {
	h.F(r)
}
