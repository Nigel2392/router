package router

import (
	"net/http"

	"github.com/Nigel2392/routevars"
)

// Route is a single route in the router
type Route struct {
	Method      string
	Path        string
	HandlerFunc HandleFunc
	middleware  []func(Handler) Handler
	children    []*Route
}

// HandleFunc registers a new route with the given path and method.
func (r *Route) HandleFunc(method, path string, handler HandleFunc) Registrar {
	path = r.Path + path
	var child = &Route{
		Method:      method,
		Path:        path,
		HandlerFunc: handler,
		middleware:  make([]func(Handler) Handler, 0),
	}
	r.children = append(r.children, child)
	return child
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
func (r *Route) Group(path string, middlewares ...func(Handler) Handler) Registrar {
	var route = &Route{Path: path}
	route.middleware = append(r.middleware, middlewares...)
	r.children = append(r.children, route)
	return route
}

// Add a group to the route
func (r *Route) AddGroup(group Registrar) {
	var g = group.(*Route)
	if g.middleware == nil {
		g.middleware = make([]func(Handler) Handler, 0)
	}
	g.middleware = append(g.middleware, r.middleware...)
	r.children = append(r.children, g)
}

// Match checks if the given path matches the route
func (r *Route) Match(method, path string) (bool, *Route, Vars) {
	if r.Method != ALL {
		if r.Method != method && r.HandlerFunc != nil {
			return false, nil, nil
		}
	}
	var ok, vars = routevars.Match(r.Path, path)
	if ok {
		return true, r, vars
	}
	for _, child := range r.children {
		if ok, route, vars := child.Match(method, path); ok {
			return ok, route, vars
		}
	}
	return false, nil, nil
}

// Use adds middleware to the route
func (r *Route) Use(middlewares ...func(Handler) Handler) {
	r.middleware = append(r.middleware, middlewares...)
}
