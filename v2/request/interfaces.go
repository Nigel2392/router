package request

// Default request user interface.
// This interface is used to check if a user is authenticated.
// This interface is used by the LoginRequiredMiddleware and LogoutRequiredMiddleware.
// If you want to use these middlewares, you should implement this interface.
// And set the GetRequestUserFunc function to return a user.
type User interface {
	IsAuthenticated() bool
}

// This interface will be set on the request, but is only useful if any middleware
// is using it. If no middleware has set it, it will remain unused.
type Session interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Exists(key string) bool
	Delete(key string)
	Destroy() error
}

// Default logger interface, can be used to set a logger on the request.
// This logger can be set in for example, the middleware, and then be used in the views by the request.
type Logger interface {
	// Write an error message, loglevel error
	Error(format any, args ...any)
	// Write a warning message, loglevel warning
	Warning(format any, args ...any)
	// Write an info message, loglevel info
	Info(format any, args ...any)
	// Write a debug message, loglevel debug
	Debug(format any, args ...any)
	// Write a test message, loglevel test
	Test(format any, args ...any)
}
