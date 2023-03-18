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

// Default json response which gets returned when using (j).Encode().
type JSONResponse struct {
	Next   string         `json:"next,omitempty"`
	Detail string         `json:"detail,omitempty"`
	Status ResponseStatus `json:"status"`
	Data   interface{}    `json:"data"`
}

// Encode json to a request.
func JSON(r *request.Request, jsonResponse *JSONResponse) error {
	var jsonData, err = json.Marshal(jsonResponse)
	if err != nil {
		return err
	}
	r.Response.Header().Set("Content-Type", "application/json")
	r.Response.Write(jsonData)
	return nil
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
func Encode(r *request.Request, data interface{}, status ...ResponseStatus) error {
	var response = JSONResponse{
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
func Decode(r *request.Request, data interface{}) error {
	// Check header
	if r.Request.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type is not application/json")
	}
	return json.NewDecoder(r.Request.Body).Decode(data)
}
