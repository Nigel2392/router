package middleware

import (
	"net/http"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

func Recoverer(next router.Handler) router.Handler {
	return router.HandleFunc(func(r *request.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(r.Response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(r)
	})
}
