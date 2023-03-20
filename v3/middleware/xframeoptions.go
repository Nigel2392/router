package middleware

import (
	"github.com/Nigel2392/router/v3"
	"github.com/Nigel2392/router/v3/request"
)

// XFrameOption is the type for the XFrameOptions middleware.
type XFrameOption string

const (
	// XFrameDeny is the most restrictive option, and it tells the browser to not display the content in an iframe.
	XFrameDeny XFrameOption = "DENY"
	// XFrameSame is the default value for XFrameOptions.
	XFrameSame XFrameOption = "SAMEORIGIN"
	// XFrameAllow is a special case, and it should not be used.
	// It is obsolete and is only here for backwards compatibility.
	XFrameAllow XFrameOption = "ALLOW-FROM"
)

// X-Frame-Options is a header that can be used to indicate whether or not a browser should be allowed to render a page in a <frame>, <iframe> or <object>.
//
// Sites can use this to avoid clickjacking attacks, by ensuring that their content is not embedded into other sites.
func XFrameOptions(options XFrameOption) router.Middleware {
	return func(next router.Handler) router.Handler {
		return router.HandleFunc(func(r *request.Request) {
			r.Response.Header().Set("X-Frame-Options", string(options))
			next.ServeHTTP(r)
		})
	}
}
