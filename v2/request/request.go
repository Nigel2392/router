package request

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
)

// Default request user interface.
// This interface is used to check if a user is authenticated.
// This interface is used by the LoginRequiredMiddleware and LogoutRequiredMiddleware.
// If you want to use these middlewares, you should implement this interface.
// And set the GetRequestUserFunc function to return a user.
type User interface {
	IsAuthenticated() bool
}

// You need to override this function to be able to use the `request.User` field.
// Function which will be called at the initialization of every request.
// This function should return a user, to set on the request.
//
// Beware!
//
// - This function is called before any middlewares are called!
var GetRequestUserFunc = func(w http.ResponseWriter, r *http.Request) User {
	return nil
}

type RequestConstraint interface {
	*Request | *http.Request
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
	Data      map[string]interface{}
	URLParams URLParams
	form      url.Values
	User      User
	JSON      *_json
}

// Initialize a new request.
func NewRequest(writer http.ResponseWriter, request *http.Request, params URLParams) *Request {
	var r = &Request{
		Response:  writer,
		Request:   request,
		URLParams: params,
		JSON:      &_json{},
	}
	r.JSON.r = &r
	r.User = GetRequestUserFunc(r.Response, r.Request)
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

// Check if the request method is POST.
func (r *Request) IsPost() bool {
	return r.Request.Method == http.MethodPost
}

// Check if the request method is GET.
func (r *Request) IsGet() bool {
	return r.Request.Method == http.MethodGet
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
	if r.Data == nil {
		r.Data = make(map[string]interface{})
	}
	r.Data[key] = value
}

// Get a data value.
func (r *Request) GetData(key string) interface{} {
	if r.Data == nil {
		return nil
	}
	return r.Data[key]
}
