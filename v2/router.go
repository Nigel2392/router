package router

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nigel2392/router/v2/request"
)

// HTTP Methods
const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"

	// Special method for matching all methods
	ALL = "ALL"
)

// Registrar is the main interface for registering routes
type Registrar interface {
	// Put registers a new route with the given path and method.
	Put(path string, handler HandleFunc) Registrar

	// Post registers a new route with the given path and method.
	Post(path string, handler HandleFunc) Registrar

	// Get registers a new route with the given path and method.
	Get(path string, handler HandleFunc) Registrar

	// Delete registers a new route with the given path and method.
	Delete(path string, handler HandleFunc) Registrar

	// Patch registers a new route with the given path and method.
	Patch(path string, handler HandleFunc) Registrar

	// Options registers a new route with the given path and method.
	Options(path string, handler HandleFunc) Registrar

	// Head registers a new route with the given path and method.
	Head(path string, handler HandleFunc) Registrar

	// Register a route for all methods
	Any(path string, handler HandleFunc) Registrar

	// HandleFunc registers a new route with the given path and method.
	HandleFunc(method, path string, handler HandleFunc) Registrar

	// Handle is a convenience method that wraps the http.Handler in a HandleFunc
	Handle(method, path string, handler http.Handler) Registrar

	// Use adds middleware to the router.
	Use(middlewares ...Middleware)

	// Group creates a new router URL group
	Group(path string, middlewares ...Middleware) Registrar

	// Addgroup adds a group of routes to the router
	AddGroup(group Registrar)

	// This is the only function the router does not implement.
	// Formats the URL for the given route, based on the given arguments.
	URL(args ...any) string
}

// Middleware is the function that is called when a route is matched
type Middleware func(Handler) Handler

// Handler is the interface that wraps the ServeHTTP method.
type Handler interface {
	ServeHTTP(*request.Request)
}

// HandleFunc is the function that is called when a route is matched
type HandleFunc func(*request.Request)

// Variable map passed to the route.
type Vars map[string]string

// Router is the main router struct
// It takes care of dispatching requests to the correct route
type Router struct {
	routes            []*Route
	middleware        []Middleware
	skipTrailingSlash bool
}

// NewRouter creates a new router
func NewRouter(SkipTrailingSlash bool) *Router {
	var r = &Router{routes: make([]*Route, 0), middleware: make([]Middleware, 0), skipTrailingSlash: SkipTrailingSlash}
	return r
}

// HandleFunc registers a new route with the given path and method.
func (r *Router) HandleFunc(method, path string, handler HandleFunc) Registrar {
	var route = &Route{Method: method, Path: path, HandlerFunc: handler, middlewareEnabled: true, middleware: r.middleware}
	r.routes = append(r.routes, route)
	return route
}

// Handle is a convenience method that wraps the http.Handler in a HandleFunc
func (r *Router) Handle(method, path string, handler http.Handler) Registrar {
	return r.HandleFunc(method, path, HTTPWrapper(handler.ServeHTTP))
}

// Put registers a new route with the given path and method.
func (r *Router) Put(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("PUT", path, handler)
}

// Post registers a new route with the given path and method.
func (r *Router) Post(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("POST", path, handler)
}

// Get registers a new route with the given path and method.
func (r *Router) Get(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("GET", path, handler)
}

// Delete registers a new route with the given path and method.
func (r *Router) Delete(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("DELETE", path, handler)
}

// Patch registers a new route with the given path and method.
func (r *Router) Patch(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("PATCH", path, handler)
}

// Options registers a new route with the given path and method.
func (r *Router) Options(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("OPTIONS", path, handler)
}

// Head registers a new route with the given path and method.
func (r *Router) Head(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("HEAD", path, handler)
}

// Register a route for all methods
func (r *Router) Any(path string, handler HandleFunc) Registrar {
	return r.HandleFunc(ALL, path, handler)
}

// Use adds middleware to the router.
func (r *Router) Use(middlewares ...Middleware) {
	r.middleware = append(r.middleware, middlewares...)
}

// Group creates a new router URL group
func (r *Router) Group(path string, middlewares ...Middleware) Registrar {
	var route = &Route{Path: path, middlewareEnabled: true}
	r.routes = append(r.routes, route)
	route.middleware = append(r.middleware, middlewares...)
	return route
}

// Addgroup adds a group of routes to the router
func (r *Router) AddGroup(group Registrar) {
	r.routes = append(r.routes, group.(*Route))
}

// ServeHTTP dispatches the request to the handler whose
// pattern matches the request URL.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if r.skipTrailingSlash && len(req.URL.Path) > 1 && req.URL.Path[len(req.URL.Path)-1] == '/' {
		req.URL.Path = req.URL.Path[:len(req.URL.Path)-1]
	}

	for _, route := range r.routes {
		if ok, newRoute, vars := route.Match(req.Method, req.URL.Path); ok && newRoute.HandlerFunc != nil {

			// Create a new handler
			var handler Handler = handleFuncWrapper{newRoute.HandlerFunc}

			// Only run the global middleware if the
			// route has middleware enabled
			if newRoute.middlewareEnabled && len(r.middleware) > 0 {
				for i := len(r.middleware) - 1; i >= 0; i-- {
					handler = r.middleware[i](handler)
				}
			}

			// Run the route middleware
			for i := len(newRoute.middleware) - 1; i >= 0; i-- {
				handler = newRoute.middleware[i](handler)
			}

			// Serve the request
			handler.ServeHTTP(request.NewRequest(w, req, vars))
			return
		}
	}
	http.NotFound(w, req)
}

var replacer = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "'", "&apos;")

// SafePath escapes the path for XML
func safePath(path string) string {
	return replacer.Replace(path)
}

// SiteMap returns a ready to use XML sitemap
func (r *Router) SiteMap() []byte {
	var maxDepth int
	for _, route := range r.routes {
		walkRoutes(route, 1, func(route *Route, depth int) {
			if depth > maxDepth {
				maxDepth = depth
			}
		})
	}

	var priority = func(depth int) float64 {
		return 1.0 - (float64(depth) / float64(maxDepth))
	}

	var buffer bytes.Buffer
	buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	buffer.WriteString("	<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	for _, route := range r.routes {
		WalkRoutes(route, func(route *Route, depth int) {
			var d = priority(depth)
			if route.HandlerFunc != nil {
				buffer.WriteString("		<url>\n")
				buffer.WriteString("			<loc>" + safePath(route.Path) + "</loc>\n")
				buffer.WriteString("			<priority>" + fmt.Sprintf("%.2f", d) + "</priority>\n")
				buffer.WriteString("		</url>\n")
			}
		})
	}
	buffer.WriteString(`	</urlset>`)

	return buffer.Bytes()
}
