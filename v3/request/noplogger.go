package request

type NopLogger struct{}

func NewNopLogger() Logger {
	return &NopLogger{}
}

func (l *NopLogger) Test(args ...interface{})                 {}
func (l *NopLogger) Testf(format string, args ...interface{}) {}

func (l *NopLogger) Debug(args ...interface{})                 {}
func (l *NopLogger) Debugf(format string, args ...interface{}) {}

func (l *NopLogger) Info(args ...interface{})                 {}
func (l *NopLogger) Infof(format string, args ...interface{}) {}

func (l *NopLogger) Warning(args ...interface{})                 {}
func (l *NopLogger) Warningf(format string, args ...interface{}) {}

func (l *NopLogger) Error(args ...interface{})                 {}
func (l *NopLogger) Errorf(format string, args ...interface{}) {}

func (l *NopLogger) Critical(err error)                           {}
func (l *NopLogger) Criticalf(format string, args ...interface{}) {}

func (l *NopLogger) LogLevel() LogLevel { return LogLevelTest }
