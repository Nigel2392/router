package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// Define methods and encodings type
type Methods string
type Encoding string

// Define a type for multipart files
type File struct {
	FileName  string    // Name of the file
	FieldName string    // Name of the field
	Reader    io.Reader // Reader of the file
}

// Define standard error messages
var (
	ErrNoRequest  = "no request has been set"                  // Error message for when no request has been set
	ErrNoCallback = "no callback has been set"                 // Error message for when no callback has been set
	ErrNoEncoding = "no encoding has been set or is not valid" // Error message for when no encoding has been set
)

// Define request methods
const (
	GET     Methods = "GET"     // GET method
	POST    Methods = "POST"    // POST method
	PUT     Methods = "PUT"     // PUT method
	PATCH   Methods = "PATCH"   // PATCH method
	DELETE  Methods = "DELETE"  // DELETE method
	OPTIONS Methods = "OPTIONS" // OPTIONS method
	HEAD    Methods = "HEAD"    // HEAD method
	TRACE   Methods = "TRACE"   // TRACE method

)

// Define methods of encoding
const (
	FORM_URL_ENCODED Encoding = "application/x-www-form-urlencoded" // FORM_URL_ENCODED encoding
	MULTIPART_FORM   Encoding = "multipart/form-data"               // MULTIPART_FORM encoding
	JSON             Encoding = "json"                              // JSON encoding
	XML              Encoding = "xml"                               // XML encoding
)

// Client is a client that can be used to make requests to a server.
func NewClient() *Client {
	return &Client{
		client: &http.Client{},
	}
}

// Initialize a GET request
func (c *Client) Request(r *http.Request) *Client {
	c.request = r
	return c
}

// Add form data to the request
func (c *Client) WithData(formData map[string]string, encoding Encoding, file ...File) error {
	if c.request == nil {
		return errors.New(ErrNoRequest)
	}

	switch encoding {
	case JSON:
		c.request.Header.Set("Content-Type", string(JSON))
		buf := new(bytes.Buffer)
		var err = json.NewEncoder(buf).Encode(formData)
		if err != nil {
			return err
		}
		c.request.Body = io.NopCloser(buf)

	case FORM_URL_ENCODED:
		c.request.Header.Set("Content-Type", string(FORM_URL_ENCODED))
		var formValues = url.Values{}
		for k, v := range formData {
			formValues.Add(k, v)
		}
		c.request.Body = io.NopCloser(bytes.NewBufferString(formValues.Encode()))

	case MULTIPART_FORM:
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for k, v := range formData {
			writer.WriteField(k, v)
		}
		for _, f := range file {
			part, err := writer.CreateFormFile(f.FieldName, f.FileName)
			if err != nil {
				return err
			}
			_, err = io.Copy(part, f.Reader)
			if err != nil {
				return err
			}
		}
		c.request.Header.Set("Content-Type", writer.FormDataContentType())
		c.request.Body = io.NopCloser(body)
	default:
		return errors.New(ErrNoEncoding)
	}
	return nil
}

// Make a request with url query parameters
func (c *Client) WithQuery(query map[string]string) error {
	if c.request == nil {
		return errors.New(ErrNoRequest)
	}
	q := c.request.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	c.request.URL.RawQuery = q.Encode()
	return nil
}

// Add headers to the request
func (c *Client) WithHeaders(headers map[string]string) error {
	if c.request == nil {
		return errors.New(ErrNoRequest)
	}
	for k, v := range headers {
		c.request.Header.Set(k, v)
	}
	return nil
}

// Add a HTTP.Cookie to the request
func (c *Client) WithCookie(cookie *http.Cookie) error {
	if c.request == nil {
		return errors.New(ErrNoRequest)
	}
	c.request.AddCookie(cookie)
	return nil
}

// Do not reccover when an error occurs
func (c *Client) OnRecover(f func(err error)) *Client {
	c.onRecover = f
	return c
}

// Recover from a panic and print the stack trace
func Recover(f func(err error)) any {
	if r := recover(); r != nil {
		if f != nil {
			f(r.(error))
		}
		return r
	}
	return nil
}
