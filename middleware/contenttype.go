package middleware

import (
	"net/http"
	"strings"

	"github.com/Nigel2392/router"
)

func AllowContentType(contentTypes ...string) func(next router.Handler) router.Handler {
	return func(next router.Handler) router.Handler {
		allowedContentTypes := make(map[string]int8, len(contentTypes))
		for _, ctype := range contentTypes {
			allowedContentTypes[strings.TrimSpace(strings.ToLower(ctype))] = 0
		}

		return router.HandleFuncWrapper{F: func(v router.Vars, w http.ResponseWriter, r *http.Request) {

			// Check if the body has content
			if r.ContentLength == 0 {
				next.ServeHTTP(v, w, r)
				return
			}

			var cTyp = r.Header.Get("Content-Type")
			var s = strings.ToLower(strings.TrimSpace(cTyp))
			if i := strings.Index(s, ";"); i > -1 {
				s = s[0:i]
			}

			if _, ok := allowedContentTypes[s]; ok {
				next.ServeHTTP(v, w, r)
				return
			}

			http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		}}
	}
}
