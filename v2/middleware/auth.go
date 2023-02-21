package middleware

import (
	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

// Middleware that only allows users who are authenticated to continue.
// By default, will call the notAuth function.
// Set the following function to change the default behavior:
//
//	request.GetRequestUserFunc = func(r *request.Request) request.User {
//		user = ...
//		return user
//	}
func LoginRequiredMiddleware(notAuth func(r *request.Request)) func(next router.Handler) router.Handler {
	if notAuth == nil {
		panic("LoginRequiredMiddleware: notAuth function is nil")
	}
	return func(next router.Handler) router.Handler {
		return router.ToHandler(func(r *request.Request) {
			if r.User == nil || !r.User.IsAuthenticated() {
				notAuth(r)
			} else {
				next.ServeHTTP(r)
			}
		})
	}
}

// Middleware that only allows users who are authenticated to continue.
// By default, will always redirect.
// Set the following function to change the default behavior:
//
//	request.GetRequestUserFunc = func(r *request.Request) request.User {
//		user = ...
//		return user
//	}
func LoginRequiredRedirectMiddleware(nextURL string) func(next router.Handler) router.Handler {
	return LoginRequiredMiddleware(func(r *request.Request) {
		router.RedirectWithNextURL(r, nextURL)
	})
}

// Middleware that only allows users who are not authenticated to continue
// By default, will never call the isAuth function.
// Set the following function to change the default behavior:
//
//	request.GetRequestUserFunc = func(r *request.Request) request.User {
//		user = ...
//		return user
//	}
func LogoutRequiredMiddleware(isAuth func(r *request.Request)) func(next router.Handler) router.Handler {
	if isAuth == nil {
		panic("LogoutRequiredMiddleware: isAuth function is nil")
	}
	return func(next router.Handler) router.Handler {
		return router.ToHandler(func(r *request.Request) {
			if r.User != nil && r.User.IsAuthenticated() {
				isAuth(r)
			} else {
				next.ServeHTTP(r)
			}
		})
	}
}

// Middleware that only allows users who are not authenticated to continue
// By default, will never call the isAuth function.
// Set the following function to change the default behavior:
//
//	request.GetRequestUserFunc = func(r *request.Request) request.User {
//		user = ...
//		return user
//	}
func LogoutRequiredRedirectMiddleware(nextURL string) func(next router.Handler) router.Handler {
	return LogoutRequiredMiddleware(func(r *request.Request) {
		router.RedirectWithNextURL(r, nextURL)
	})
}
