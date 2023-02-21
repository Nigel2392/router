package middleware

import (
	"fmt"
	"net/http"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

// Check if the request.Host is in the allowed hosts list
func AllowedHosts(allowed_hosts ...string) func(next router.Handler) router.Handler {
	if len(allowed_hosts) == 0 {
		panic("AllowedHosts: No hosts provided.")
	}
	for _, host := range allowed_hosts {
		if host == "" {
			panic("AllowedHosts: Empty host not allowed.")
		} else if host == "*" {
			// If the host is set to *, allow all hosts
			return func(next router.Handler) router.Handler {
				return next
			}
		}
	}
	return func(next router.Handler) router.Handler {
		return router.ToHandler(func(r *request.Request) {
			// Check if ALLOWED_HOSTS is set and if the request host is allowed
			var allowed = false
			var requestHost = request.GetHost(r)
			for _, host := range allowed_hosts {
				if host == requestHost {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(r.Writer, fmt.Sprintf("Host not allowed: %s", requestHost), http.StatusForbidden)
				return
			}
			next.ServeHTTP(r)
		})
	}
}
