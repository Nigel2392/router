package middleware

import (
	"github.com/Nigel2392/router/v2/request"
)

var DEFAULT_LOGGER request.Logger

// Format the message, paired with the request IP and method.
func formatMessage(r *request.Request, format string, args ...any) string {
	return logFormat(r, "MIDDLEWARE ERROR", format, args...)
}
