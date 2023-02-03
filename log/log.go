package log

import (
	"github.com/jingyanbin/core/internal"
)

func SetLevel(level LOGGER_LEVEL) {
	log.SetLevel(level)
}

func Debug(v ...any) {
	if log.Level() > DEBUG {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.output("D", file, line, v...)
}

func Info(v ...any) {
	if log.Level() > INFO {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.output("I", file, line, v...)
}

func Warn(v ...any) {
	if log.Level() > WARN {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.output("W", file, line, v...)
}

func Error(v ...any) {
	if log.Level() > ERROR {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.output("E", file, line, v...)
}

func Fatal(v ...any) {
	if log.Level() > FATAL {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.output("F", file, line, v...)
}

func DebugF(format string, v ...any) {
	if log.Level() > DEBUG {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.outputF("D", file, line, format, v...)
}

func InfoF(format string, v ...any) {
	if log.Level() > INFO {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.outputF("I", file, line, format, v...)
}

func WarnF(format string, v ...any) {
	if log.Level() > WARN {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.outputF("W", file, line, format, v...)
}

func ErrorF(format string, v ...any) {
	if log.Level() > ERROR {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.outputF("E", file, line, format, v...)
}

func FatalF(format string, v ...any) {
	if log.Level() > FATAL {
		return
	}
	file, line := internal.CallerShort(logSkip)
	log.outputF("F", file, line, format, v...)
}

func Close() {
	log.Close()
}
