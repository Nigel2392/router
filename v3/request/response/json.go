package response

import (
	"encoding/json"
	"errors"

	"github.com/Nigel2392/router/v3/request"
)

type ResponseStatus string

const (
	ResponseStatusOK       ResponseStatus = "ok"
	ResponseStatusError    ResponseStatus = "error"
	ResponseStatusRedirect ResponseStatus = "redirect"
)

type JSONResponse struct {
	Next   string         `json:"next,omitempty"`
	Detail string         `json:"detail,omitempty"`
	Status ResponseStatus `json:"status"`
	Data   interface{}    `json:"data"`
}

// Encode json to a request.
func Json(r *request.Request, jsonResponse *JSONResponse) error {
	var jsonData, err = json.Marshal(jsonResponse)
	if err != nil {
		return err
	}
	r.Response.Header().Set("Content-Type", "application/json")
	r.Response.Write(jsonData)
	return nil
}

// Render a JSONError to a request.
type JSONError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"` // HTTP status code

	// The error that caused this error
	Err error `json:"-"`
}

func (e *JSONError) Error() string {
	return e.Message
}

func (e *JSONError) Unwrap() error {
	return e.Err
}

// Create a new JSONError
func NewJsonError(message string, statusCode int, err error) *JSONError {
	return &JSONError{
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

func (e *JSONError) Write(r *request.Request) error {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return err
	}
	r.Response.Header().Set("Content-Type", "application/json")
	r.Response.Write(jsonData)
	return nil
}

func JsonError(r *request.Request, msg string, statusCode int, err error) error {
	var response = JSONError{
		Message:    msg,
		StatusCode: statusCode,
		Err:        err,
	}

	return response.Write(r)
}

// Render json to a request.
// Response will be in the form of:
//
//	{
//		"status": "ok",
//		"data": {
//			"key": "value"
//		}
//	}
func JsonEncode(r *request.Request, data interface{}, status ...ResponseStatus) error {
	var response = JSONResponse{
		Next: r.Next(),
		Data: data,
	}
	if len(status) > 0 {
		response.Status = status[0]
	} else {
		response.Status = ResponseStatusOK
	}
	var jsonData, err = json.Marshal(response)
	if err != nil {
		return err
	}
	r.Response.Header().Set("Content-Type", "application/json")
	r.Response.Write(jsonData)
	return nil
}

// Decoode json from a request, into any.
func JsonDecode(r *request.Request, data interface{}) error {
	// Check header
	if r.Request.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type is not application/json")
	}
	return json.NewDecoder(r.Request.Body).Decode(data)
}
