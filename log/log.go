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

func Close() {
	log.Close()
}
