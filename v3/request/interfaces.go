package request

// Default request user interface.
//
// This interface is used to check if a user is authenticated.
//
// This interface is used by the LoginRequiredMiddleware and LogoutRequiredMiddleware.
//
// If you want to use these middlewares, you should implement this interface.
//
// And set the GetRequestUserFunc function to return a user.
type User interface {
	// Check if the user is authenticated
	IsAuthenticated() bool

	// Check if the user is an administator
	IsAdmin() bool

	// Check permissions, if the user is an admin, this should return true.
	HasPermissions(permissions ...string) bool
}

// This interface will be set on the request, but is only useful if any middleware
//
// is using it. If no middleware has set it, it will remain unused.
type Session interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Exists(key string) bool
	Delete(key string)
	Destroy() error
	RenewToken() error
}

type LogLevel int

const (
	LogLevelCritical LogLevel = iota + 1
	LogLevelError
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
	LogLevelTest
)

// Default logger interface, can be used to set a logger on the request.
//
// This logger can be set in for example, the middleware, and then be used in the views by the request.
type Logger interface {
	// Write a critical error message
	//
	// This message should be handled differently
	//
	// than the other ways of reporting.
	Critical(err error)
	// Format and write a critical error message
	Criticalf(format string, args ...any)
	// Write an error message, loglevel error
	Error(args ...any)
	Errorf(format string, args ...any)
	// Write a warning message, loglevel warning
	Warning(args ...any)
	Warningf(format string, args ...any)
	// Write an info message, loglevel info
	Info(args ...any)
	Infof(format string, args ...any)
	// Write a debug message, loglevel debug
	Debug(args ...any)
	Debugf(format string, args ...any)
	// Write a test message, loglevel test
	Test(args ...any)
	Testf(format string, args ...any)
	// Retrieve the loglevel
	LogLevel() LogLevel
}
