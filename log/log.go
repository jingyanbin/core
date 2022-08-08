package log

import (
	"github.com/jingyanbin/core/datetime"
	"github.com/jingyanbin/core/internal"
)

//import . "github.com/jingyanbin/basal"
//import . "github.com/jingyanbin/datetime"

func SetLevel(level int) {
	log.SetLevel(level)
}

func AddHandler(handler logWriter) {
	log.AddHandler(handler)
}

func SetAsync(async bool) {
	log.SetAsync(async)
}

func Debug(v ...interface{}) {
	log.output(LOG_LEVEL_DEBUG, v...)
}

func Info(v ...interface{}) {
	log.output(LOG_LEVEL_INFO, v...)
}

func Warn(v ...interface{}) {
	log.output(LOG_LEVEL_WARN, v...)
}

func Error(v ...interface{}) {
	log.output(LOG_LEVEL_ERROR, v...)
}

func Fatal(v ...interface{}) {
	log.output(LOG_LEVEL_FATAL, v...)
}

func DebugF(format string, v ...interface{}) {
	log.outputf(LOG_LEVEL_DEBUG, format, v...)
}

func InfoF(format string, v ...interface{}) {
	log.outputf(LOG_LEVEL_INFO, format, v...)
}

func WarnF(format string, v ...interface{}) {
	log.outputf(LOG_LEVEL_WARN, format, v...)
}

func ErrorF(format string, v ...interface{}) {
	log.outputf(LOG_LEVEL_ERROR, format, v...)
}

func FatalF(format string, v ...interface{}) {
	log.outputf(LOG_LEVEL_FATAL, format, v...)
}

func SetFormatHeader(formatHeader func(buf *internal.Buffer, level string, line int, file string, dt *datetime.DateTime)) {
	log.SetFormatHeader(formatHeader)
}

func Wait() { //等待日志模块退出
	loggerMgr.Wait()
	log.Wait()
}
