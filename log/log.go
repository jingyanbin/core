package log

import (
	"github.com/jingyanbin/core/internal"
)

type LOGGER_LEVEL = internal.LOGGER_LEVEL

const (
	DEBUG = internal.DEBUG
	INFO  = internal.INFO
	WARN  = internal.WARN
	ERROR = internal.ERROR
	FATAL = internal.FATAL
	OFF   = internal.OFF
)

func SetILogger(logger internal.ILogger) {
	internal.Log = logger
}

func GetILogger() internal.ILogger {
	return internal.Log
}

func SetLevel(level LOGGER_LEVEL) {
	internal.StdLog.SetLevel(level)
}

func Debug(v ...any) {
	if internal.StdLog.Level() > DEBUG {
		return
	}
	file, line := internal.CallerShort(internal.LogSkip)
	internal.StdLog.Output("D", file, line, v...)
}

func Info(v ...any) {
	if internal.StdLog.Level() > INFO {
		return
	}
	file, line := internal.CallerShort(internal.LogSkip)
	internal.StdLog.Output("I", file, line, v...)
}

func Warn(v ...any) {
	if internal.StdLog.Level() > WARN {
		return
	}
	file, line := internal.CallerShort(internal.LogSkip)
	internal.StdLog.Output("W", file, line, v...)
}

func Error(v ...any) {
	if internal.StdLog.Level() > ERROR {
		return
	}
	file, line := internal.CallerShort(internal.LogSkip)
	internal.StdLog.Output("E", file, line, v...)
}

func Fatal(v ...any) {
	if internal.StdLog.Level() > FATAL {
		return
	}
	file, line := internal.CallerShort(internal.LogSkip)
	internal.StdLog.Output("F", file, line, v...)
}

func Close() {
	internal.StdLog.Close()
}
