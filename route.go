package router

import (
	"net/http"

	"github.com/Nigel2392/routevars"
)

type Route struct {
	Method      string
	Path        string
	HandlerFunc HandleFunc
	middleware  []func(Handler) Handler
	children    []*Route
}

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

func (r *Route) Handle(method, path string, handler http.Handler) Registrar {
	return r.HandleFunc(method, path, HTTPWrapper(handler.ServeHTTP))
}

func (r *Route) Put(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("PUT", path, handler)
}

func (r *Route) Post(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("POST", path, handler)
}

func (r *Route) Get(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("GET", path, handler)
}

func (r *Route) Delete(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("DELETE", path, handler)
}

func (r *Route) Patch(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("PATCH", path, handler)
}

func (r *Route) Options(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("OPTIONS", path, handler)
}

func (r *Route) Head(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("HEAD", path, handler)
}

func (r *Route) Match(method, path string) (bool, Vars) {
	if r.Method != method {
		return false, nil
	}
	var ok, vars = routevars.Match(r.Path, path)
	if ok {
		return true, vars
	}
	for _, child := range r.children {
		if ok, vars = child.Match(method, path); ok {
			return ok, vars
		}
	}
	return false, nil
}

func (r *Route) Use(middlewares ...func(Handler) Handler) {
	r.middleware = append(r.middleware, middlewares...)
}
