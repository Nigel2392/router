package middleware

import (
	"fmt"

	"github.com/Nigel2392/router/v2/request"
)

var Logger request.Logger

// Format the message, paired with the request IP and method.
func formatMessage(r *request.Request, format string, args ...any) string {
	return fmt.Sprintf("%s %s %s", r.IP().String(), r.Method(), fmt.Sprintf(format, args...))
}
