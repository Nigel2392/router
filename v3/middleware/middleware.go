package middleware

import (
	"fmt"
	"time"

	"github.com/Nigel2392/router/v3/request"
)

var DEFAULT_LOGGER request.Logger

// Format the message, paired with the request IP and method.
func FormatMessage(r *request.Request, messageType string, format string, args ...any) string {
	return fmt.Sprintf("[%s %s %s] %s %s %s",
		r.IP(),
		r.Method(),
		time.Now().Format("2006-01-02 15:04:05"),
		messageType,
		r.Request.URL.Path,
		fmt.Sprintf(format, args...))
}
