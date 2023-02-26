package middleware

import (
	"net/http"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

// Recoverer recovers from panics and logs the error,
// if the logger was set.
func Recoverer(next router.Handler) router.Handler {
	return router.HandleFunc(func(r *request.Request) {
		defer func() {
			if err := recover(); err != nil {
				if DEFAULT_LOGGER != nil {
					DEFAULT_LOGGER.Error(formatMessage(r, "Panic: %s", err))
				}
				http.Error(r.Response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(r)
	})
}
