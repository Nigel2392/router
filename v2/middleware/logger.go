package middleware

import (
	"fmt"
	"runtime"
	"strings"
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

var (
	critMsg  = "\u001B[31;1;4mCRITICAL\u001B[0m"
	errMsg   = "   \u001B[31;4mERROR\u001B[0m"
	warnMsg  = " \u001B[33;4mWARNING\u001B[0m"
	infoMsg  = "    \u001B[34;4mINFO\u001B[0m"
	dbugMsg  = "   \u001B[32;1;4mDEBUG\u001B[0m"
	testMsg  = "    \u001B[35;1;4mTEST\u001B[0m"
	writeMSG = "   \u001B[90;4mWRITE\u001B[0m"
)

// Log a critical message.
func (l *logger) Critical(err error) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("[ \u001B[90;4m%s - \u001B[90m%s\u001B[0m %s ] %s on line %d \u001B[90m%s\u001B[0m\n",
		padString(l.request.Method()+"\u001B[0m", 11),
		time.Now().Format("2006-01-02 15:04:05"),
		critMsg,
		err.Error(),
		line,
		file)
}

// Error logs an error message.
func (l *logger) Error(args ...any) {
	logPrintln(l.request, errMsg, args...)
}

// Warning logs a warning message.
func (l *logger) Warning(args ...any) {
	logPrintln(l.request, warnMsg, args...)
}

// Info logs an info message.
func (l *logger) Info(args ...any) {
	logPrintln(l.request, infoMsg, args...)
}

// Debug logs a debug message.
func (l *logger) Debug(args ...any) {
	logPrintln(l.request, dbugMsg, args...)
}

// Test logs a test message.
func (l *logger) Test(args ...any) {
	logPrintln(l.request, testMsg, args...)
}

// Error logs an error message.
func (l *logger) Errorf(format string, args ...any) {
	logPrintf(l.request, errMsg, format, args...)
}

// Warning logs a warning message.
func (l *logger) Warningf(format string, args ...any) {
	logPrintf(l.request, warnMsg, format, args...)
}

// Info logs an info message.
func (l *logger) Infof(format string, args ...any) {
	logPrintf(l.request, infoMsg, format, args...)
}

// Debug logs a debug message.
func (l *logger) Debugf(format string, args ...any) {
	logPrintf(l.request, dbugMsg, format, args...)
}

// Test logs a test message.
func (l *logger) Testf(format string, args ...any) {
	logPrintf(l.request, testMsg, format, args...)
}

// Write a message to the console.
func (l *logger) Write(p []byte) (n int, err error) {
	fmt.Print(logFormat(l.request, writeMSG, string(p)))
	return len(p), nil
}

// Format a message and print it to the console.
func logPrintf(r *request.Request, levelMessage, format string, args ...any) {
	fmt.Println(logFormat(r, levelMessage, fmt.Sprintf(format, args...)))
}

// Format a message and print it to the console.
func logPrintln(r *request.Request, levelMessage string, args ...any) {
	fmt.Print(logFormat(r, levelMessage, fmt.Sprintln(args...)))
}

// Format a message and return it.
func logFormat(r *request.Request, levelMessage, additional string) string {
	return fmt.Sprintf("[ \u001B[90;4m%s - \u001B[90m%s\u001B[0m %s ] \u001B[90m%s\u001B[0m %s",
		padString(r.Method()+"\u001B[0m", 11),
		time.Now().Format("2006-01-02 15:04:05"),
		levelMessage,
		r.Request.URL.Path,
		additional)
}

func padString(s string, length int) string {
	var b strings.Builder
	b.Grow(length)
	b.WriteString(s)
	for i := len(s); i < length; i++ {
		b.WriteRune(' ')
	}
	return b.String()
}
