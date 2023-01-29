package router

import (
	"net/http"
)

type HandleFunc func(Vars, http.ResponseWriter, *http.Request)

type HandleFuncWrapper struct {
	F func(Vars, http.ResponseWriter, *http.Request)
}

func (h HandleFuncWrapper) ServeHTTP(vars Vars, w http.ResponseWriter, req *http.Request) {
	h.F(vars, w, req)
}

type Registrar interface {
	Put(path string, handler HandleFunc) Registrar
	Post(path string, handler HandleFunc) Registrar
	Get(path string, handler HandleFunc) Registrar
	Delete(path string, handler HandleFunc) Registrar
	Patch(path string, handler HandleFunc) Registrar
	Options(path string, handler HandleFunc) Registrar
	Head(path string, handler HandleFunc) Registrar
	HandleFunc(method, path string, handler HandleFunc) Registrar
	Handle(method, path string, handler http.Handler) Registrar
	Use(middlewares ...func(Handler) Handler)
	Group(path string, middlewares ...func(Handler) Handler) Registrar
}

type Handler interface {
	ServeHTTP(Vars, http.ResponseWriter, *http.Request)
}

type Vars map[string]string

func HTTPWrapper(handler func(http.ResponseWriter, *http.Request)) HandleFunc {
	return func(vars Vars, w http.ResponseWriter, req *http.Request) {
		handler(w, req)
	}
}

type Router struct {
	routes     []*Route
	middleware []func(Handler) Handler
	http.Handler
	skipTrailingSlash bool
}

func NewRouter(SkipTrailingSlash bool) *Router {
	var r = &Router{routes: make([]*Route, 0), middleware: make([]func(Handler) Handler, 0), skipTrailingSlash: SkipTrailingSlash}
	return r
}

func (r *Router) HandleFunc(method, path string, handler HandleFunc) Registrar {
	var route = &Route{Method: method, Path: path, HandlerFunc: handler}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) Handle(method, path string, handler http.Handler) Registrar {
	return r.HandleFunc(method, path, HTTPWrapper(handler.ServeHTTP))
}

func (r *Router) Put(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("PUT", path, handler)
}

func (r *Router) Post(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("POST", path, handler)
}

func (r *Router) Get(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("GET", path, handler)
}

func (r *Router) Delete(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("DELETE", path, handler)
}

func (r *Router) Patch(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("PATCH", path, handler)
}

func (r *Router) Options(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("OPTIONS", path, handler)
}

func (r *Router) Head(path string, handler HandleFunc) Registrar {
	return r.HandleFunc("HEAD", path, handler)
}

func (r *Router) Use(middlewares ...func(Handler) Handler) {
	r.middleware = append(r.middleware, middlewares...)
}

func (r *Router) Group(path string, middlewares ...func(Handler) Handler) Registrar {
	var route = &Route{Path: path}
	route.middleware = append(r.middleware, middlewares...)
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if r.skipTrailingSlash && len(req.URL.Path) > 1 && req.URL.Path[len(req.URL.Path)-1] == '/' {
		req.URL.Path = req.URL.Path[:len(req.URL.Path)-1]
	}

	for _, route := range r.routes {
		if ok, newRoute, vars := route.Match(req.Method, req.URL.Path); ok && newRoute.HandlerFunc != nil {

			var handler Handler = HandleFuncWrapper{newRoute.HandlerFunc}

			for _, middleware := range r.middleware {
				handler = middleware(handler)
			}

			for _, middleware := range newRoute.middleware {
				handler = middleware(handler)
			}

			handler.ServeHTTP(vars, w, req)
			return
		}
	}
	http.NotFound(w, req)
}
