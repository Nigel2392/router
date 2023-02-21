package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Nigel2392/router/v2/request"
	"github.com/Nigel2392/routevars"
)

// Route is a single route in the router
type Route struct {
	Method            string
	Path              string
	HandlerFunc       HandleFunc
	middleware        []Middleware
	children          []*Route
	middlewareEnabled bool
}

// HandleFunc registers a new route with the given path and method.
func (r *Route) HandleFunc(method, path string, handler HandleFunc) Registrar {
	path = r.Path + path
	var child = &Route{
		Method:            method,
		Path:              path,
		HandlerFunc:       handler,
		middleware:        make([]Middleware, 0),
		middlewareEnabled: true,
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
func (r *Route) Put(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("PUT", path, handler)
}

// Post registers a new route with the given path and method.
func (r *Route) Post(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("POST", path, handler)
}

// Get registers a new route with the given path and method.
func (r *Route) Get(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("GET", path, handler)
}

// Delete registers a new route with the given path and method.
func (r *Route) Delete(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("DELETE", path, handler)
}

// Patch registers a new route with the given path and method.
func (r *Route) Patch(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("PATCH", path, handler)
}

// Options registers a new route with the given path and method.
func (r *Route) Options(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("OPTIONS", path, handler)
}

// Head registers a new route with the given path and method.
func (r *Route) Head(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("HEAD", path, handler)
}

// Register a route for all methods
func (r *Route) Any(path string, handler HandleFunc) Registrar {
	return r.HandleFunc(ALL, path, handler)
}

// Group creates a new group of routes
func (r *Route) Group(path string, middlewares ...Middleware) Registrar {
	var route = &Route{
		Path:       r.Path + path,
		middleware: middlewares,
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
	for _, child := range g.children {
		child.Path = r.Path + child.Path
	}
	r.children = append([]*Route{g}, r.children...)
}

// Match checks if the given path matches the route
func (r *Route) Match(method, path string) (bool, *Route, request.URLParams) {
	if r.Method != ALL {
		if r.Method != method && r.HandlerFunc != nil {
			return false, nil, nil
		}
	}
	if r.HandlerFunc != nil {
		var ok, vars = routevars.Match(r.Path, path)
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

// Format the url based on the arguments given.
// Panics if route accepts more arguments than are given.
func (r *Route) URL(args ...any) string {
	var path = r.Path

	// If the length of the path is less than the length of the pre/suffix and the delimiter
	// then there are no variables in the path
	if len(path) <= len(routevars.RT_PATH_VAR_DELIM)+len(routevars.RT_PATH_VAR_PREFIX)+len(routevars.RT_PATH_VAR_SUFFIX) {
		return path
	}
	// Remove the first and last slash if they exist
	var hasPrefixSlash = strings.HasPrefix(path, "/")
	var hasTrailingSlash = strings.HasSuffix(path, "/")
	if hasPrefixSlash {
		path = path[1:]
	}
	if hasTrailingSlash {
		path = path[:len(path)-1]
	}
	// Split the path into parts
	var parts = strings.Split(path, "/")
	// Replace the parts that are variables with the arguments
	for i, part := range parts {
		if strings.HasPrefix(part, routevars.RT_PATH_VAR_PREFIX) && strings.HasSuffix(part, routevars.RT_PATH_VAR_SUFFIX) {
			if len(args) == 0 {
				panic("not enough arguments for URL: " + r.Path)
			}
			var arg = args[0]
			args = args[1:]
			parts[i] = fmt.Sprintf("%v", arg)
		}
	}
	// Join the parts back together
	path = strings.Join(parts, "/")
	// Add the slashes back if they were there
	if hasPrefixSlash {
		path = "/" + path
	}
	if hasTrailingSlash {
		path = path + "/"
	}
	return path
}
