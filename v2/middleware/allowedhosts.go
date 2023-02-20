package middleware

import (
	"fmt"
	"net/http"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

func AllowedHosts(allowed_hosts ...string) func(next router.Handler) router.Handler {
	return func(next router.Handler) router.Handler {
		return router.HandleFuncWrapper{F: func(r *request.Request) {
			// Check if ALLOWED_HOSTS is set and if the request host is allowed
			if len(allowed_hosts) > 0 {
				var allowed = false
				var requestHost = request.GetHost(r)
				for _, host := range allowed_hosts {
					if host == requestHost || host == "*" {
						allowed = true
						break
					}
				}
				if !allowed {
					http.Error(r.Writer, fmt.Sprintf("Host not allowed: %s", requestHost), http.StatusForbidden)
					return
				}
			} else {
				// If ALLOWED_HOSTS is not set, deny all requests
				http.Error(r.Writer, "Allowed Hosts not set.", http.StatusForbidden)
			}
			next.ServeHTTP(r)
		}}
	}
}
