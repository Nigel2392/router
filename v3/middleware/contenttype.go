package middleware

import (
	"net/http"
	"strings"

	"github.com/Nigel2392/router/v3"
	"github.com/Nigel2392/router/v3/request"
)

// AllowContentType is a middleware that checks the content type of the request
// body. If the content type is not allowed, the middleware will return a
// http.StatusUnsupportedMediaType error.
func AllowContentType(contentTypes ...string) func(next router.Handler) router.Handler {
	return func(next router.Handler) router.Handler {
		allowedContentTypes := make(map[string]int8, len(contentTypes))
		for _, ctype := range contentTypes {
			allowedContentTypes[strings.TrimSpace(strings.ToLower(ctype))] = 0
		}

		return router.HandleFunc(func(r *request.Request) {

			// Check if the body has content
			if r.Request.ContentLength == 0 {
				next.ServeHTTP(r)
				return
			}

			var cTyp = r.GetHeader("Content-Type")
			var s = strings.ToLower(strings.TrimSpace(cTyp))
			if i := strings.Index(s, ";"); i > -1 {
				s = s[0:i]
			}

			if _, ok := allowedContentTypes[s]; ok {
				next.ServeHTTP(r)
				return
			}

			http.Error(r.Response, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		})
	}
}
