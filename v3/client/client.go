package client

import (
	"errors"
	"net/http"
)

// Client is a client that can be used to execute http requests.
// - Can be used to execute GET, POST, PUT, DELETE, PATCH requests.
type Client struct {
	// Client is the http client that will be used to execute the request.
	client *http.Client
	// Request is the request that will be executed.
	request *http.Request
	// Recover from panics?
	onRecover func(err error)
}

// Execute the request -> APIClient.exec
func (c *Client) Do() (*http.Response, error) {
	if c.request == nil {
		return nil, errors.New(ErrNoRequest)
	}
	var resp, err = c.client.Do(c.request)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
