package request

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const MESSAGE_COOKIE_NAME = "messages"
const NEXT_COOKIE_NAME = "next"

// Default request user interface.
// This interface is used to check if a user is authenticated.
// This interface is used by the LoginRequiredMiddleware and LogoutRequiredMiddleware.
// If you want to use these middlewares, you should implement this interface.
// And set the GetRequestUserFunc function to return a user.
type User interface {
	IsAuthenticated() bool
}

// This interface is used to retrieve the request host.
type RequestConstraint interface {
	*Request | *http.Request
}

// This interface will be set on the request, but is only useful if any middleware
// is using it. If no middleware has set it, it will remain unused.
type Session interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Exists(key string) bool
	Delete(key string)
	Destroy() error
}

func GetHost[T RequestConstraint](r T) string {
	var host string
	switch r := any(r).(type) {
	case *Request:
		host = r.Request.Host
	case *http.Request:
		host = r.Host
	}
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	return host
}

// Default request to be passed around the router.
type Request struct {
	Response  http.ResponseWriter
	Request   *http.Request
	Data      *TemplateData
	Session   Session
	URLParams URLParams
	form      url.Values
	User      User
	JSON      *_json
	next      string
}

// Initialize a new request.
func NewRequest(writer http.ResponseWriter, request *http.Request, params URLParams) *Request {
	var r = &Request{
		Response:  writer,
		Request:   request,
		URLParams: params,
		JSON:      &_json{},
		Data:      NewTemplateData(),
	}
	r.JSON.r = &r
	return r
}

// Write to the response.
func (r *Request) Write(b []byte) (int, error) {
	return r.Response.Write(b)
}

// Write a string to the response.
func (r *Request) WriteString(s string) (int, error) {
	return r.Response.Write([]byte(s))
}

// Raise an error.
func (r *Request) Error(code int, err string) {
	http.Error(r.Response, err, code)
}

// Get the request method.
func (r *Request) Method() string {
	return r.Request.Method
}

// Parse the form, and return the form values.
func (r *Request) Form() url.Values {
	if r.form == nil {
		r.Request.ParseForm()
		r.form = r.Request.Form
	}
	return r.form
}

// Get a form file as a buffer.
func (r *Request) FormFileBuffer(name string) (*bytes.Buffer, error) {
	m, _, err := r.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.ReadFrom(m)
	return &buf, nil
}

// Set a data value.
func (r *Request) SetData(key string, value interface{}) {
	r.Data.Set(key, value)
}

// Get a data value.
func (r *Request) GetData(key string) interface{} {
	return r.Data.Get(key)
}

// Get the Next url.
// This is the url that was set in the session/cookies.
// This is used to redirect back to the same page.
func (r *Request) Next() string {
	if r.next != "" {
		return r.next
	}
	if r.Session == nil {
		// Set the next url if it exists.
		// This is based on cookies.
		if cookie, err := r.Request.Cookie(NEXT_COOKIE_NAME); err == nil {
			r.next = cookie.Value
			// Delete the cookie.
			http.SetCookie(r.Response, &http.Cookie{
				Name:    NEXT_COOKIE_NAME,
				Value:   "",
				Expires: time.Now().Add(-time.Hour),
			})
		}
	} else {
		// We have sessions! :)
		if next, ok := r.Session.Get(NEXT_COOKIE_NAME).(string); ok {
			r.next = next
			r.Session.Delete(NEXT_COOKIE_NAME)
		}
	}
	return r.next
}

// Redirect to a URL.
// If the session is defined, the messages will be set in the session.
// If the `next` argument is given, it will be added to session, unless
// the session is not defined, the `next` parameter will be added to cookies.
// This means they will be carried across when rendering with request.Render().
// This will be set again after the redirect, when rendering in the default Data.
// Optionally you could obtain this by calling request.Next().
func (r *Request) Redirect(redirectURL string, statuscode int, next ...string) {
	// Set the messages in the session/cookies for after the redirect.
	if r.Session == nil {
		// If there is a next parameter, add it to the cookies.
		if len(next) > 0 && next[0] != "" {
			var cookie = &http.Cookie{
				Name:     NEXT_COOKIE_NAME,
				Value:    next[0],
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour * 24 * 30),
				Secure:   r.Request.TLS != nil,
				MaxAge:   60 * 60 * 24 * 30,
			}
			http.SetCookie(r.Response, cookie)
		}
		// Set the messages in the cookies.
		if r.Data != nil {
			var cookie = &http.Cookie{
				Name:     MESSAGE_COOKIE_NAME,
				Value:    r.Data.Messages.Encode(),
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour * 24 * 30),
				Secure:   r.Request.TLS != nil,
				MaxAge:   60 * 60 * 24 * 30,
			}
			http.SetCookie(r.Response, cookie)
		}

	} else {
		// We have sessions! :)
		if r.Data != nil {
			r.Session.Set(MESSAGE_COOKIE_NAME, r.Data.Messages)
		}
		if r.next != "" {
			r.Session.Set(NEXT_COOKIE_NAME, r.next)
		}
	}

	// Redirect.
	http.Redirect(r.Response, r.Request, redirectURL, statuscode)
}
