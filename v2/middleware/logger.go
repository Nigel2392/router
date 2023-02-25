package middleware

import (
	"fmt"
	"time"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

func AddLogger(next router.Handler) router.Handler {
	return router.HandleFunc(func(r *request.Request) {
		Logger := &logger{request: r}
		r.Logger = Logger
		next.ServeHTTP(r)
	})
}

type logger struct {
	request *request.Request
}

// Log an error message.
func (l *logger) Error(format string, args ...any) {
	fmt.Printf("%s \u001B[4m%s\u001B[0m - [%s \u001B[31;4mError\u001B[0m] %s\n",
		l.request.Method(),
		l.request.Request.URL.Path,
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(format, args...))
}

// Log a warning message.
func (l *logger) Warning(format string, args ...any) {
	fmt.Printf("%s \u001B[4m%s\u001B[0m - [%s \u001B[33;4mWarning\u001B[0m] %s\n",
		l.request.Method(),
		l.request.Request.URL.Path,
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(format, args...))
}

// Log an info message.
func (l *logger) Info(format string, args ...any) {
	fmt.Printf("%s \u001B[4m%s\u001B[0m - [%s \u001B[34;4mInfo\u001B[0m] %s\n",
		l.request.Method(),
		l.request.Request.URL.Path,
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(format, args...))
}

// Log a debug message.
func (l *logger) Debug(format string, args ...any) {
	fmt.Printf("%s \u001B[4m%s\u001B[0m - [%s \u001B[32;1;4mDebug\u001B[0m] %s\n",
		l.request.Method(),
		l.request.Request.URL.Path,
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(format, args...))
}

// Log a test message.
func (l *logger) Test(format string, args ...any) {
	fmt.Printf("%s \u001B[4m%s\u001B[0m - [%s \u001B[35;1;4mTest\u001B[0m] %s\n",
		l.request.Method(),
		l.request.Request.URL.Path,
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(format, args...))
}

// Format the message, paired with the request IP and method.
func formatMessage(r *request.Request, format string, args ...any) string {
	return fmt.Sprintf("%s %s %s", r.IP().String(), r.Method(), fmt.Sprintf(format, args...))
}
