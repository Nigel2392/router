package router

import (
	"net/http"
	"net/url"

	"github.com/Nigel2392/router/v2/request"
)

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
