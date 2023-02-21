package router

import (
	"net/http"
	"net/url"

	"github.com/Nigel2392/router/v2/request"
)

// Group creates a new router URL group
func Group(path string) Registrar {
	var route = &Route{Path: path}
	return route
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
	http.Redirect(r.Writer, r.Request, new_login_url.String(), http.StatusFound)
}

// Wrapper function for http.Handler to make it compatible with HandleFunc
func HTTPWrapper(handler func(http.ResponseWriter, *http.Request)) HandleFunc {
	return func(r *request.Request) {
		handler(r.Writer, r.Request)
	}
}

// Wrapper function for router.Handler to make it compatible with http.handler
func HandlerWrapper(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(request.NewRequest(w, req, nil))
	})
}
