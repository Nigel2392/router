package request

import (
	"net/http"
	"strings"
)

// This constraint is used to retrieve the request host.
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

// Add a header onto the request.
// This will append the value to the current value.
// If the value already exists, it will not be added.
func AddHeader(w http.ResponseWriter, name, value string) {
	// Get the current value.
	current := w.Header().Get(name)
	// If the current value is empty, set the value.
	if current == "" {
		w.Header().Set(name, value)
		return
	}
	// If the value already exists, do nothing.
	if strings.Contains(current, value) {
		return
	}
	// If the current value is not empty, append the value.
	w.Header().Set(name, current+"; "+value)
}
