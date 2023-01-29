package middleware

import (
	"net/http"

	"github.com/Nigel2392/router"
)

func Recoverer(next router.Handler) router.Handler {
	return router.HandleFuncWrapper{F: func(v router.Vars, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(v, w, r)
	}}
}
