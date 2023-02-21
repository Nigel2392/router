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
	Status ResponseStatus `json:"status"`
	Data   interface{}    `json:"data"`
}

// Intermediate struct for json encoding/decoding.
type _json struct {
	r **Request
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
func (j *_json) Encode(status ResponseStatus, data interface{}) error {
	var response = JSONResponse{
		Status: status,
		Data:   data,
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
