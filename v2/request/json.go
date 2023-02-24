package request

import (
	"encoding/json"
	"errors"
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

// Intermediate struct for json encoding/decoding.
type _json struct {
	r **Request
}

// Encode json to a request.
func (j *_json) SendResponse(jsonResponse *JSONResponse) error {
	var jsonData, err = json.Marshal(jsonResponse)
	if err != nil {
		return err
	}
	(*j.r).Response.Header().Set("Content-Type", "application/json")
	(*j.r).Response.Write(jsonData)
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
func (j *_json) Encode(data interface{}, status ...ResponseStatus) error {
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
	(*j.r).Response.Header().Set("Content-Type", "application/json")
	(*j.r).Response.Write(jsonData)
	return nil
}

// Decoode json from a request, into any.
func (j *_json) Decode(data interface{}) error {
	// Check header
	var r = (*j.r).Request
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type is not application/json")
	}
	var err = json.NewDecoder((*j.r).Request.Body).Decode(data)
	return err
}
