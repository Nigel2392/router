package middleware

import (
	"errors"
	"fmt"

	"github.com/Nigel2392/router/v3"
	"github.com/Nigel2392/router/v3/request"
)

// Recoverer recovers from panics and logs the error,
// if the logger was set.
func Recoverer(onError func(err error, r *request.Request)) router.Middleware {
	return func(next router.Handler) router.Handler {
		return router.HandleFunc(func(r *request.Request) {
			defer func() {
				if err := recover(); err != nil {
					if DEFAULT_LOGGER != nil {
						DEFAULT_LOGGER.Error(FormatMessage(r, "PANIC", "Panic: %v", err))
					}
					r.Response.Clear()
					var newHandler = router.HandleFunc(func(r *request.Request) {
						switch err := err.(type) {
						case error:
							onError(err, r)
						case string:
							onError(errors.New(err), r)
						default:
							onError(fmt.Errorf("%v", err), r)
						}
					})
					newHandler.ServeHTTP(r)
				}
			}()
			next.ServeHTTP(r)
		})
	}
}
