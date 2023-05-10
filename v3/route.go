package router

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nigel2392/router/v3/client"
	"github.com/Nigel2392/router/v3/request"
	"github.com/Nigel2392/router/v3/request/params"
	"github.com/Nigel2392/router/v3/request/writer"
	"github.com/Nigel2392/routevars"
)

// Route is a single route in the router
type Route struct {
	Method            string
	Path              routevars.URLFormatter
	HandlerFunc       Handler
	middleware        []Middleware
	children          []*Route
	middlewareEnabled bool
	name              string
}

// Return the name of the route
// This can be used to match a route with the URL method
func (r *Route) Name() string {
	return r.name
}

// Returns the formatted URL for this route.
//
// If no arguments are provided it will return the path as it is set on the route.
// This means it will be returned as /my/blog/<<post_id:int>>
func (r *Route) Format(args ...any) string {
	if len(args) <= 0 {
		return string(r.Path)
	}
	return r.Path.Format(args)
}

// URL matches a URL for the given names delimited by a colon.
// Example:
//
// Parent:Child:GrandChild
func (r *Route) URL(method, name string) routevars.URLFormatter {
	var parts = strings.Split(name, ":")
	if r.name == parts[0] && len(parts) == 1 {
		if r.Method == method ||
			r.Method == ALL ||
			method == ALL {
			return r.Path
		}
	} else if r.name == parts[0] && len(parts) > 1 {
		if r := r.url(method, parts[1:]); r != "" {
			return r
		}
	}
	return ""
}

// Route returns the route that matches the given method and path
func (r *Route) url(method string, parts []string) routevars.URLFormatter {
	if len(parts) == 0 {
		return ""
	}
	var thismatch = r.name == parts[0]
	if thismatch && r.Method == method || thismatch && r.Method == ALL {
		return ""
	}
	for _, route := range r.children {
		if len(parts) == 1 {
			var rmatch = route.name == parts[0]
			if rmatch && route.Method == method ||
				rmatch && route.Method == ALL ||
				rmatch && method == ALL {
				return route.Path
			}
		} else if len(parts) > 1 && route.name == parts[0] {
			if r := route.url(method, parts[1:]); r != "" {
				return r
			}
		}
	}
	return ""
}

// HandleFunc registers a new route with the given path and method.
func (r *Route) HandleFunc(method, path string, handler Handler, name ...string) Registrar {
	var n = r.name
	if len(name) > 0 {
		n = name[0]
	}
	path = string(r.Path) + path
	var child = &Route{
		Method:            method,
		Path:              routevars.URLFormatter(path),
		HandlerFunc:       handler,
		middleware:        r.middleware,
		middlewareEnabled: r.middlewareEnabled,
		name:              n,
	}
	r.children = append(r.children, child)
	return child
}

// Disable the router's middlewares for this route, and all its children
// It will however still run the route's own middlewares.
func (r *Route) DisableMiddleware() {
	r.middlewareEnabled = false
	for _, child := range r.children {
		child.DisableMiddleware()
	}
}

// Handle is a convenience method that wraps the http.Handler in a HandleFunc
func (r *Route) Handle(method, path string, handler http.Handler) Registrar {
	return r.HandleFunc(method, path, HTTPWrapper(handler.ServeHTTP))
}

// Put registers a new route with the given path and method.
func (r *Route) Put(path string, handler Handler, name ...string) Registrar {
	return r.HandleFunc("PUT", path, handler, name...)
}

// Post registers a new route with the given path and method.
func (r *Route) Post(path string, handler Handler, name ...string) Registrar {
	return r.HandleFunc("POST", path, handler, name...)
}

// Get registers a new route with the given path and method.
func (r *Route) Get(path string, handler Handler, name ...string) Registrar {
	return r.HandleFunc("GET", path, handler, name...)
}

// Delete registers a new route with the given path and method.
func (r *Route) Delete(path string, handler Handler, name ...string) Registrar {
	return r.HandleFunc("DELETE", path, handler, name...)
}

// Patch registers a new route with the given path and method.
func (r *Route) Patch(path string, handler Handler, name ...string) Registrar {
	return r.HandleFunc("PATCH", path, handler, name...)
}

// Options registers a new route with the given path and method.
func (r *Route) Options(path string, handler Handler, name ...string) Registrar {
	return r.HandleFunc("OPTIONS", path, handler, name...)
}

// Head registers a new route with the given path and method.
func (r *Route) Head(path string, handler Handler, name ...string) Registrar {
	return r.HandleFunc("HEAD", path, handler, name...)
}

// Register a route for all methods
func (r *Route) Any(path string, handler Handler, name ...string) Registrar {
	return r.HandleFunc(ALL, path, handler, name...)
}

// Group creates a new group of routes
func (r *Route) Group(path string, name string, middlewares ...Middleware) Registrar {
	if len(middlewares) == 0 {
		middlewares = make([]Middleware, 0)
	}
	middlewares = append(middlewares, r.middleware...)
	var route = &Route{
		Path:              r.Path + routevars.URLFormatter(path),
		middleware:        middlewares,
		middlewareEnabled: r.middlewareEnabled,
		name:              name,
	}
	r.children = append([]*Route{route}, r.children...)
	return route
}

// Add a group to the route
func (r *Route) AddGroup(group Registrar) {
	var g = group.(*Route)
	if g.middleware == nil {
		g.middleware = make([]Middleware, 0)
	}
	g.middleware = append(g.middleware, r.middleware...)
	g.Path = r.Path + g.Path
	g.middlewareEnabled = r.middlewareEnabled
	for _, child := range g.children {
		WalkRoutes(child, func(route *Route, i int) {
			route.Path = r.Path + route.Path
			route.middleware = append(route.middleware, g.middleware...)
		})
	}
	r.children = append([]*Route{g}, r.children...)
}

// Match checks if the given path matches the route
func (r *Route) Match(method, path string) (bool, *Route, params.URLParams) {
	if r.Method != ALL && r.HandlerFunc == nil {
		if r.Method != method && r.HandlerFunc != nil {
			return false, nil, nil
		}
	}
	if r.HandlerFunc != nil {
		var ok, vars = r.Path.Match(path)
		if ok {
			return true, r, vars
		}
	}
	for _, child := range r.children {
		if ok, route, vars := child.Match(method, path); ok {
			return ok, route, vars
		}
	}
	return false, nil, nil
}

// Use adds middleware to the route
func (r *Route) Use(middlewares ...Middleware) {
	r.middleware = append(r.middleware, middlewares...)
	for _, child := range r.children {
		child.Use(middlewares...)
	}
}

// Call a route handler with the given request.
//
// Do so by making a HTTP request to the route's url.
//
// If te url takes arguments, you need to pass them into Call.
//
// It will run the route's middleware and the route's handler.
//
// This will
func (r *Route) Call(req *http.Request, args ...any) (*http.Response, error) {
	var handler = r.HandlerFunc
	var path = r.Path.Format(args...)
	if handler == nil {
		panic(fmt.Sprintf("No handler for route %s", path))
	}
	req.URL.Path = path
	var client = client.NewClient()
	client.OnRecover(func(err error) {})
	client.Request(req)
	return client.Do()
}

// Invoke a route handler directly.
//
// If te url takes arguments, you need to pass them into Invoke.
//
// It will not run the route's middleware.
func (r *Route) Invoke(dest http.ResponseWriter, req *http.Request, args ...any) {
	var handler = r.HandlerFunc
	var path = r.Path.Format(args...)
	if handler == nil {
		panic(fmt.Sprintf("No handler for route %s", path))
	}
	req.URL.Path = path
	var _, vars = r.Path.Match(path)
	var resp = &writer.ClearableBufferedResponseWriter{
		ResponseWriter: dest,
		Buf:            new(bytes.Buffer),
	}
	var request = request.NewRequest(resp, req, vars)
	defer resp.Finalize()
	handler.ServeHTTP(request)
}
