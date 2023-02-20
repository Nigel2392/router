package request

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
)

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

type Request struct {
	Writer    http.ResponseWriter
	Request   *http.Request
	Data      map[string]interface{}
	URLParams URLParams
	form      url.Values
	JSON      *_json
}

func NewRequest(writer http.ResponseWriter, request *http.Request, params URLParams) *Request {
	var r = &Request{
		Writer:    writer,
		Request:   request,
		URLParams: params,
	}
	r.JSON.r = &r
	return r
}

func (r *Request) Error(code int, err string) {
	http.Error(r.Writer, err, code)
}

func (r *Request) Method() string {
	return r.Request.Method
}

func (r *Request) IsPost() bool {
	return r.Request.Method == http.MethodPost
}

func (r *Request) IsGet() bool {
	return r.Request.Method == http.MethodGet
}

func (r *Request) Form() url.Values {
	if r.form == nil {
		r.Request.ParseForm()
		r.form = r.Request.Form
	}
	return r.form
}

func (r *Request) FormFileBuffer(name string) (*bytes.Buffer, error) {
	m, _, err := r.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.ReadFrom(m)
	return &buf, nil
}

func (r *Request) SetData(key string, value interface{}) {
	if r.Data == nil {
		r.Data = make(map[string]interface{})
	}
	r.Data[key] = value
}

func (r *Request) GetData(key string) interface{} {
	if r.Data == nil {
		return nil
	}
	return r.Data[key]
}
