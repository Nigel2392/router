package csrf

import (
	"encoding/base64"
	"net/http"

	"github.com/Nigel2392/router/v2/request"
)

/*

PACKAGE BASED ON nosurf!

https://github.com/justinas/nosurf

________________________________________________________________________________

The MIT License (MIT)

Copyright (c) 2013 Justinas Stankevicius

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// Default errors to use
const (
	ErrTokenNotFound       = "CSRF token not found"
	ErrTokenMismatch       = "CSRF token mismatch"
	ErrRefererNotSpecified = "Referer not specified"
	ErrRefererMismatch     = "Referer mismatch"
)

// List of unsafe HTTP methods
var unsafeMethods = listCompare[string]{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

const (
	tokenLength = 32
)

// Get the token for a given request once the middleware has been run
// If the token is not found, an empty string is returned
func Token(req *request.Request) string {
	val, _ := contextGetToken(req)
	return val
}

// VerifyToken verifies the sent token equals the real one
// and returns a bool value indicating if tokens are equal.
// Supports masked tokens. realToken comes from Token(r) and
// sentToken is token sent unusual way.
func VerifyToken(realToken, sentToken string) bool {
	r, err := base64.StdEncoding.DecodeString(realToken)
	if err != nil {
		return false
	}
	if len(r) == 2*tokenLength {
		r = unmaskToken(r)
	}
	s, err := base64.StdEncoding.DecodeString(sentToken)
	if err != nil {
		return false
	}
	if len(s) == 2*tokenLength {
		s = unmaskToken(s)
	}
	return tokensEqual(r, s)
}

// verifyToken expects the realToken to be unmasked and the sentToken to be masked
func verifyToken(realToken, sentToken []byte) bool {
	realN := len(realToken)
	sentN := len(sentToken)

	// sentN == tokenLength means the token is unmasked
	// sentN == 2*tokenLength means the token is masked.

	if realN == tokenLength && sentN == 2*tokenLength {
		return tokensEqual(realToken, unmaskToken(sentToken))
	}
	return false
}

// Extracts the "sent" token from the request
// and returns an unmasked version of it
func extractToken(r *request.Request) []byte {
	// Prefer the header over form value
	sentToken := r.Request.Header.Get(CSRF_TOKEN_HEADER_NAME)

	// Then POST values
	if len(sentToken) == 0 {
		sentToken = r.Request.PostFormValue(CSRF_TOKEN_FORMFIELD_NAME)
	}

	// If all else fails, try a multipart value.
	// PostFormValue() will already have called ParseMultipartForm()
	if len(sentToken) == 0 && r.Request.MultipartForm != nil {
		vals := r.Request.MultipartForm.Value[CSRF_TOKEN_FORMFIELD_NAME]
		if len(vals) != 0 {
			sentToken = vals[0]
		}
	}

	return b64decode(sentToken)
}
