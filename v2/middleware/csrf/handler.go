package csrf

import (
	"net/http"
	"net/url"
	"time"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

const (
	CSRF_TOKEN_COOKIE_NAME      = "csrf_token"
	CSRF_TOKEN_HEADER_NAME      = "X-CSRF-Token"
	CSRF_TOKEN_FORMFIELD_NAME   = "csrf_token"
	CSRF_TOKEN_COOKIE_EXPIRE    = time.Hour * 24
	CSRF_TOKEN_COOKIE_MAX_AGE   = 3600
	CSRF_TOKEN_COOKIE_SECURE    = false
	CSRF_TOKEN_COOKIE_HTTP_ONLY = true
)

func Middleware(next router.Handler) router.Handler {
	return router.HandleFunc(func(req *request.Request) {

		defaultContext(req)

		req.AddHeader("Vary", "Cookie")

		var realToken []byte
		tokenCookie, err := req.GetCookie(CSRF_TOKEN_COOKIE_NAME)
		if err == nil {
			realToken = b64decode(tokenCookie.Value)
		}

		if len(realToken) != tokenLength {
			var t = generateToken()
			contextSaveToken(req, b64encode(maskToken(realToken)))
			var cookie = &http.Cookie{
				Name:     CSRF_TOKEN_COOKIE_NAME,
				Value:    b64encode(t),
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
			}
			req.SetCookies(cookie)
		} else {
			contextSaveToken(req, b64encode(maskToken(realToken)))
		}

		if req.Data == nil {
			req.Data = request.NewTemplateData()
		}

		req.Data.CSRFToken = request.NewCSRFToken(Token(req))

		// Check if the request method is safe.
		if !unsafeMethods.Contains(req.Method()) {
			// Continue to the next handler.
			next.ServeHTTP(req)
			return
		}

		if req.Request.URL.Scheme == "https" {
			referer, err := url.Parse(req.GetHeader("Referer"))

			// if we can't parse the referer or it's empty,
			// we assume it's not specified
			if err != nil || referer.String() == "" {
				// Error: Referer not specified
				req.Error(http.StatusForbidden, ErrRefererNotSpecified)
				return
			}

			// if the referer doesn't share origin with the request URL,
			// we have another error for that
			if referer.Scheme != req.Request.URL.Scheme || referer.Host != req.Request.URL.Host {
				// Error: Referer mismatch
				req.Error(http.StatusForbidden, ErrRefererMismatch)
				return
			}
		}

		// Finally, we check the token itself.
		sentToken := extractToken(req)

		if !verifyToken(realToken, sentToken) {
			// Error: Token mismatch
			req.Error(http.StatusForbidden, ErrTokenMismatch)
			return
		}

		// Continue to the next handler.
		next.ServeHTTP(req)
	})
}
