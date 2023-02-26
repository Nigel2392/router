package middleware

import (
	"fmt"
	"time"

	"github.com/Nigel2392/router/v2"
	"github.com/Nigel2392/router/v2/request"
)

func AddLogger(next router.Handler) router.Handler {
	return router.HandleFunc(func(r *request.Request) {
		r.Logger = &logger{request: r}
		next.ServeHTTP(r)
	})
}

type logger struct {
	request *request.Request
}

// Error logs an error message.
func (l *logger) Error(args ...any) {
	logPrintln(l.request, "\u001B[31;4mERROR\u001B[0m  ", args...)
}

// Warning logs a warning message.
func (l *logger) Warning(args ...any) {
	logPrintln(l.request, "\u001B[33;4mWARNING\u001B[0m", args...)
}

// Info logs an info message.
func (l *logger) Info(args ...any) {
	logPrintln(l.request, "\u001B[34;4mINFO\u001B[0m   ", args...)
}

// Debug logs a debug message.
func (l *logger) Debug(args ...any) {
	logPrintln(l.request, "\u001B[32;1;4mDEBUG\u001B[0m  ", args...)
}

// Test logs a test message.
func (l *logger) Test(args ...any) {
	logPrintln(l.request, "\u001B[35;1;4mTEST\u001B[0m ", args...)
}

// Error logs an error message.
func (l *logger) Errorf(format string, args ...any) {
	logPrintf(l.request, "\u001B[31;4mERROR\u001B[0m  ", format, args...)
}

// Warning logs a warning message.
func (l *logger) Warningf(format string, args ...any) {
	logPrintf(l.request, "\u001B[33;4mWARNING\u001B[0m", format, args...)
}

// Info logs an info message.
func (l *logger) Infof(format string, args ...any) {
	logPrintf(l.request, "\u001B[34;4mINFO\u001B[0m   ", format, args...)
}

// Debug logs a debug message.
func (l *logger) Debugf(format string, args ...any) {
	logPrintf(l.request, "\u001B[32;1;4mDEBUG\u001B[0m  ", format, args...)
}

// Test logs a test message.
func (l *logger) Testf(format string, args ...any) {
	logPrintf(l.request, "\u001B[35;1;4mTEST\u001B[0m ", format, args...)
}

// Write a message to the console.
func (l *logger) Write(p []byte) (n int, err error) {
	fmt.Print(logFormat(l.request, "\u001B[90;4mWRITE\u001B[0m  ", string(p)))
	return len(p), nil
}

// Format a message and print it to the console.
func logPrintf(r *request.Request, levelMessage, format string, args ...any) {
	fmt.Println(logFormat(r, levelMessage, fmt.Sprintf(format, args...)))
}

// Format a message and print it to the console.
func logPrintln(r *request.Request, levelMessage string, args ...any) {
	fmt.Println(logFormat(r, levelMessage, fmt.Sprint(args...)))
}

// Format a message and return it.
func logFormat(r *request.Request, levelMessage, format string) string {
	return fmt.Sprintf("[\u001B[90;4m%s\u001B[0m - \u001B[90m%s\u001B[0m %s] \u001B[90m%s\u001B[0m %s",
		r.Method(),
		time.Now().Format("2006-01-02 15:04:05"),
		levelMessage,
		r.Request.URL.Path,
		format)
}
