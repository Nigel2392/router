package router

import (
	"io/fs"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nigel2392/router/v3/request"
	"github.com/Nigel2392/router/v3/templates"
	"github.com/Nigel2392/routevars"
)

// Group creates a new router URL group
func Group(path string, name string) Registrar {
	var route = &Route{Path: routevars.URLFormatter(path), middlewareEnabled: true}
	route.name = name
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

func NewFSRoute(static_url string, name string, fs fs.FS) Registrar {
	var newStatic = templates.NicePath(false, "/", static_url, "/<<any>>")
	var static_route = Group(
		newStatic, "static",
	)
	var r = static_route.(*Route)
	r.Method = GET
	r.HandlerFunc = FromHTTPHandler(http.StripPrefix(
		wrapSlash(static_url),
		http.FileServer(
			http.FS(fs),
		),
	)).ServeHTTP
	return r
}

func wrapSlash(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	if !strings.HasSuffix(p, "/") {
		p = p + "/"
	}
	return p
}
