package router

import (
	"net/http"
	"net/url"

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

// Make a new middleware function from a http.Handler
func FromHTTPMiddleware(h http.Handler) Middleware {
	return func(next Handler) Handler {
		return httpHandlerWrapper{H: h}
	}
}

type httpHandlerWrapper struct {
	H http.Handler
}

func (h httpHandlerWrapper) ServeHTTP(r *request.Request) {
	h.H.ServeHTTP(r.Response, r.Request)
}

// Group creates a new router URL group
func Group(path string) Registrar {
	var route = &Route{Path: path}
	return route
}

// Recurse over a routes children, keeping track of depth
func WalkRoutes(route *Route, f func(*Route, int)) {
	walkRoutes(route, 0, f)
}

// Recurse over a routes children, keeping track of depth
func walkRoutes(route *Route, depth int, f func(*Route, int)) {
	f(route, depth)
	for _, child := range route.children {
		walkRoutes(child, depth+1, f)
	}
}

// Redirect user to a URL, appending the current URL as a "next" query parameter
func RedirectWithNextURL(r *request.Request, nextURL string) {
	var u = r.Request.URL.String()
	var new_login_url, err = url.Parse(nextURL)
	if err != nil {
		panic(err)
	}
	var query = new_login_url.Query()
	query.Set("next", u)
	new_login_url.RawQuery = query.Encode()
	http.Redirect(r.Response, r.Request, new_login_url.String(), http.StatusFound)
}

// Wrapper function for http.Handler to make it compatible with HandleFunc
func HTTPWrapper(handler func(http.ResponseWriter, *http.Request)) HandleFunc {
	return func(r *request.Request) {
		handler(r.Response, r.Request)
	}
}

// Wrapper function for router.Handler to make it compatible with http.handler
func HandlerWrapper(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(request.NewRequest(w, req, nil))
	})
}
