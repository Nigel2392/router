package router

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nigel2392/router/v3/request"
	"github.com/Nigel2392/router/v3/request/params"
	"github.com/Nigel2392/router/v3/request/writer"
	"github.com/Nigel2392/routevars"
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
	Put(path string, handler HandleFunc, name ...string) Registrar

	// Post registers a new route with the given path and method.
	Post(path string, handler HandleFunc, name ...string) Registrar

	// Get registers a new route with the given path and method.
	Get(path string, handler HandleFunc, name ...string) Registrar

	// Delete registers a new route with the given path and method.
	Delete(path string, handler HandleFunc, name ...string) Registrar

	// Patch registers a new route with the given path and method.
	Patch(path string, handler HandleFunc, name ...string) Registrar

	// Options registers a new route with the given path and method.
	Options(path string, handler HandleFunc, name ...string) Registrar

	// Head registers a new route with the given path and method.
	Head(path string, handler HandleFunc, name ...string) Registrar

	// Register a route for all methods
	Any(path string, handler HandleFunc, name ...string) Registrar

	// HandleFunc registers a new route with the given path and method.
	HandleFunc(method, path string, handler HandleFunc, name ...string) Registrar

	// Handle is a convenience method that wraps the http.Handler in a HandleFunc
	Handle(method, path string, handler http.Handler) Registrar

	// Use adds middleware to the router.
	Use(middlewares ...Middleware)

	// Group creates a new router URL group
	Group(path string, name string, middlewares ...Middleware) Registrar

	// Addgroup adds a group of routes to the router
	AddGroup(group Registrar)

	// URL returns the URL for a named route
	URL(method, name string) routevars.URLFormatter

	// Call the route, returning the response and a possible error.
	Call(request *http.Request, args ...any) (*http.Response, error)

	// Invoke the route's handler, writing to a response writer.
	Invoke(dest http.ResponseWriter, req *http.Request, args ...any)
}

// Variable map passed to the route.
type Vars map[string]string

// Router is the main router struct
// It takes care of dispatching requests to the correct route
type Router struct {
	NotFoundHandler   Handler
	routes            []*Route
	middleware        []Middleware
	skipTrailingSlash bool
}

// Returns all the routes in a nicely formatted string for debugging.
func (r *Router) String() string {
	var buf bytes.Buffer
	for _, route := range r.routes {
		walkRoutes(route, 0, func(r *Route, i int) {
			if strings.TrimSpace(r.Method) == "" {
				fmt.Fprintf(&buf, "%s%s -> %s\n", strings.Repeat("  ", i), string(r.Path), r.name)
			} else {
				fmt.Fprintf(&buf, "%s%s %s -> %s\n", strings.Repeat("  ", i), r.Method, string(r.Path), r.name)
			}
		})
	}
	return buf.String()
}

// NewRouter creates a new router
func NewRouter(skipTrailingSlash bool) *Router {
	var r = &Router{routes: make([]*Route, 0), middleware: make([]Middleware, 0), skipTrailingSlash: skipTrailingSlash}
	return r
}

// Get a route by name.
//
// Route names are optional, when used a route's child can be access like so:
//
// It looks like Django's URL routing syntax.
//
// router.Route("routeName")
//
// router.Route("parentName:routeToGet")
func (r *Router) URL(method, name string) routevars.URLFormatter {
	var parts = strings.Split(name, ":")
	for _, route := range r.routes {
		if len(parts) == 0 {
			return ""
		}
		if route.name == parts[0] && len(parts) == 1 {
			if route.Method == method ||
				route.Method == ALL ||
				method == ALL {
				return route.Path
			}
		} else if route.name == parts[0] && len(parts) > 1 {
			if r := route.url(method, parts[1:]); r != "" {
				return r
			}
		}
	}
	return ""
}

// The URL func, but for easy use in templates.
//
// It returns the URL, formatted based on the arguments.
func (r *Router) URLFormat(name string, args ...interface{}) string {
	var url = r.URL(ALL, name)
	return url.Format(args...)
}

// HandleFunc registers a new route with the given path and method.
func (r *Router) HandleFunc(method, path string, handler HandleFunc, name ...string) Registrar {
	var route = &Route{Method: method, Path: routevars.URLFormatter(path), HandlerFunc: handler, middlewareEnabled: true}

	if len(name) > 0 {
		route.name = name[0]
	}

	r.routes = append(r.routes, route)
	return route
}

// Handle is a convenience method that wraps the http.Handler in a HandleFunc
func (r *Router) Handle(method, path string, handler http.Handler) Registrar {
	return r.HandleFunc(method, path, HTTPWrapper(handler.ServeHTTP))
}

// Put registers a new route with the given path and method.
func (r *Router) Put(path string, handler HandleFunc, name ...string) Registrar {
	return r.HandleFunc("PUT", path, handler, name...)
}

// Post registers a new route with the given path and method.
func (r *Router) Post(path string, handler HandleFunc, name ...string) Registrar {
	return r.HandleFunc("POST", path, handler, name...)
}

// Get registers a new route with the given path and method.
func (r *Router) Get(path string, handler HandleFunc, name ...string) Registrar {
	return r.HandleFunc("GET", path, handler, name...)
}

// Delete registers a new route with the given path and method.
func (r *Router) Delete(path string, handler HandleFunc, name ...string) Registrar {
	return r.HandleFunc("DELETE", path, handler, name...)
}

// Patch registers a new route with the given path and method.
func (r *Router) Patch(path string, handler HandleFunc, name ...string) Registrar {
	return r.HandleFunc("PATCH", path, handler, name...)
}

// Options registers a new route with the given path and method.
func (r *Router) Options(path string, handler HandleFunc, name ...string) Registrar {
	return r.HandleFunc("OPTIONS", path, handler, name...)
}

// Head registers a new route with the given path and method.
func (r *Router) Head(path string, handler HandleFunc, name ...string) Registrar {
	return r.HandleFunc("HEAD", path, handler, name...)
}

// Register a route for all methods
func (r *Router) Any(path string, handler HandleFunc, name ...string) Registrar {
	return r.HandleFunc(ALL, path, handler, name...)
}

// Use adds middleware to the router.
func (r *Router) Use(middlewares ...Middleware) {
	r.middleware = append(r.middleware, middlewares...)
}

// Group creates a new router URL group
func (r *Router) Group(path string, name string, middlewares ...Middleware) Registrar {
	var route = &Route{Path: routevars.URLFormatter(path), middlewareEnabled: true, name: name}
	r.routes = append(r.routes, route)
	return route
}

// Addgroup adds a group of routes to the router
func (r *Router) AddGroup(group Registrar) {
	r.routes = append(r.routes, group.(*Route))
}

// Match returns the route that matches the given method and path.
func (r *Router) Match(method, path string) (bool, *Route, params.URLParams) {
	for _, route := range r.routes {
		if ok, newRoute, vars := route.Match(method, path); ok {
			return ok, newRoute, vars
		}
	}
	return false, nil, nil
}

// ServeHTTP dispatches the request to the handler whose
// pattern matches the request URL.
func (r *Router) ServeHTTP(w http.ResponseWriter, rq *http.Request) {

	if r.skipTrailingSlash && len(rq.URL.Path) > 1 && rq.URL.Path[len(rq.URL.Path)-1] == '/' {
		rq.URL.Path = rq.URL.Path[:len(rq.URL.Path)-1]
	}

	var ok, newRoute, vars = r.Match(rq.Method, rq.URL.Path)
	if !ok {
		if r.NotFoundHandler != nil {
			var resp = writer.NewClearable(w)
			defer resp.Finalize()
			r.NotFoundHandler.ServeHTTP(request.NewRequest(resp, rq, nil))
			return
		}
		http.NotFound(w, rq)
		return
	}

	// Create a new handler
	var handler Handler = newRoute.HandlerFunc

	// Run the route middleware
	for i := len(newRoute.middleware) - 1; i >= 0; i-- {
		handler = newRoute.middleware[i](handler)
	}

	// Only run the global middleware if the
	// route has middleware enabled
	if newRoute.middlewareEnabled && len(r.middleware) > 0 {
		for i := len(r.middleware) - 1; i >= 0; i-- {
			handler = r.middleware[i](handler)
		}
	}

	// Initialize a new request.
	var req = request.NewRequest(writer.NewClearable(w), rq, vars)

	// Defer the response finalization
	//
	// This is done to actually write to the response, instead of
	// just buffering it.
	defer req.Response.Finalize()

	// Set up a function to fetch routes, from any path inside a request.
	req.URL = r.URL

	// Serve the request
	handler.ServeHTTP(req)
}

//	var replacer = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "'", "&apos;")
//
//	// SafePath escapes the path for XML
//	func safePath(path string) string {
//		return replacer.Replace(path)
//	}

//	// SiteMap returns a ready to use XML sitemap
//	func (r *Router) SiteMap() []byte {
//		var maxDepth int
//		for _, route := range r.routes {
//			walkRoutes(route, 1, func(route *Route, depth int) {
//				if depth > maxDepth {
//					maxDepth = depth
//				}
//			})
//		}
//
//		var priority = func(depth int) float64 {
//			return 1.0 - (float64(depth) / float64(maxDepth))
//		}
//
//		var buffer bytes.Buffer
//		buffer.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
//		buffer.WriteString("	<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
//		for _, route := range r.routes {
//			WalkRoutes(route, func(route *Route, depth int) {
//				var d = priority(depth)
//				if route.HandlerFunc != nil {
//					buffer.WriteString("		<url>\n")
//					buffer.WriteString("			<loc>" + safePath(string(route.Path)) + "</loc>\n")
//					buffer.WriteString("			<priority>" + fmt.Sprintf("%.2f", d) + "</priority>\n")
//					buffer.WriteString("		</url>\n")
//				}
//			})
//		}
//		buffer.WriteString(`	</urlset>`)
//
//		return buffer.Bytes()
//	}
