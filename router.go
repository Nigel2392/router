package router

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// HandleFunc is the function that is called when a route is matched
type HandleFunc func(Vars, http.ResponseWriter, *http.Request)

// Wrapper function for HandleFunc to make it compatible with http.Handler
type HandleFuncWrapper struct {
	F func(Vars, http.ResponseWriter, *http.Request)
}

// ServeHTTP implements the Handler interface
func (h HandleFuncWrapper) ServeHTTP(vars Vars, w http.ResponseWriter, req *http.Request) {
	h.F(vars, w, req)
}

// Registrar is the main interface for registering routes
// Both the router and the route struct implement this interface
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
	Any(path string, handler HandleFunc) Registrar
	Use(middlewares ...func(Handler) Handler)
	Group(path string, middlewares ...func(Handler) Handler) Registrar
}

// Handler is the interface that wraps the ServeHTTP method.
type Handler interface {
	ServeHTTP(Vars, http.ResponseWriter, *http.Request)
}

// Variable map passed to the route.
type Vars map[string]string

// Wrapper function for http.Handler to make it compatible with HandleFunc
func HTTPWrapper(handler func(http.ResponseWriter, *http.Request)) HandleFunc {
	return func(vars Vars, w http.ResponseWriter, req *http.Request) {
		handler(w, req)
	}
}

// Router is the main router struct
// It takes care of dispatching requests to the correct route
type Router struct {
	routes     []*Route
	middleware []func(Handler) Handler
	http.Handler
	skipTrailingSlash bool
}

// NewRouter creates a new router
func NewRouter(SkipTrailingSlash bool) *Router {
	var r = &Router{routes: make([]*Route, 0), middleware: make([]func(Handler) Handler, 0), skipTrailingSlash: SkipTrailingSlash}
	return r
}

// HandleFunc registers a new route with the given path and method.
func (r *Router) HandleFunc(method, path string, handler HandleFunc) Registrar {
	var route = &Route{Method: method, Path: path, HandlerFunc: handler}
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
func (r *Router) Use(middlewares ...func(Handler) Handler) {
	r.middleware = append(r.middleware, middlewares...)
}

// Group creates new children to a route.
func (r *Router) Group(path string, middlewares ...func(Handler) Handler) Registrar {
	var route = &Route{Path: path}
	route.middleware = append(r.middleware, middlewares...)
	r.routes = append(r.routes, route)
	return route
}

// ServeHTTP dispatches the request to the handler whose
// pattern matches the request URL.
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

// SiteMap returns a ready to use XML sitemap
func (r *Router) SiteMap() []byte {
	var maxDepth int
	for _, route := range r.routes {
		WalkRoutes(route, 1, func(route *Route, depth int) {
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
		WalkRoutes(route, 0, func(route *Route, depth int) {
			var d = priority(depth)
			fmt.Println(route.Path, d)
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

// Recurse over a routes children, keeping track of depth
func WalkRoutes(route *Route, depth int, f func(*Route, int)) {
	f(route, depth)
	for _, child := range route.children {
		WalkRoutes(child, depth+1, f)
	}
}

var replacer = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "'", "&apos;")

// SafePath escapes the path for XML
func safePath(path string) string {
	return replacer.Replace(path)
}

type RobotListing struct {
	Allow      []string
	Disallow   []string
	UserAgent  string
	CrawlDelay int
}

type RobotsOptions struct {
	Rules   []*RobotListing
	SiteMap string
}

// Robots returns a ready to use robots.txt
func (r *Router) Robots(options *RobotsOptions) []byte {
	var buffer bytes.Buffer
	for i, listing := range options.Rules {
		if listing.UserAgent == "" {
			listing.UserAgent = "*"
		}
		buffer.WriteString("User-agent: " + listing.UserAgent + "\n")
		for _, allow := range listing.Allow {
			buffer.WriteString("Allow: " + allow + "\n")
		}
		for _, disallow := range listing.Disallow {
			buffer.WriteString("Disallow: " + disallow + "\n")
		}
		if listing.CrawlDelay > 0 {
			buffer.WriteString("Crawl-delay: " + strconv.Itoa(listing.CrawlDelay) + "\n")
		}
		if i < len(options.Rules)-1 {
			buffer.WriteString("\n")
		}
	}
	if options.SiteMap != "" {
		buffer.WriteString("\nSitemap: " + options.SiteMap + "\n")
	}
	return buffer.Bytes()
}
