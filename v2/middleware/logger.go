package middleware

import (
	"fmt"

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

func (l *logger) Error(format string, args ...any) {
	fmt.Printf("[\u001B[31mError\u001B[0m]: %s\n", fmt.Sprintf(format, args...))
}

func (l *logger) Warning(format string, args ...any) {
	fmt.Printf("[\u001B[33mWarning\u001B[0m]: %s\n", fmt.Sprintf(format, args...))
}

func (l *logger) Info(format string, args ...any) {
	fmt.Printf("[\u001B[Info\u001B[0m]: %s\n", fmt.Sprintf(format, args...))
}

func (l *logger) Debug(format string, args ...any) {
	fmt.Printf("[\u001B[Debug\u001B[0m]: %s\n", fmt.Sprintf(format, args...))
}

func (l *logger) Test(format string, args ...any) {
	fmt.Printf("[\u001B[Test\u001B[0m]: %s\n", fmt.Sprintf(format, args...))
}

// Format the message, paired with the request IP and method.
func formatMessage(r *request.Request, format string, args ...any) string {
	return fmt.Sprintf("%s %s %s", r.IP().String(), r.Method(), fmt.Sprintf(format, args...))
}
