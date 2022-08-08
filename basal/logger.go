package basal

import (
	"github.com/jingyanbin/core/internal"
)

type ILogger = internal.ILogger

var log ILogger = GetStdoutLogger()

func GetStdoutLogger() ILogger

func SetLogger(logger ILogger) {
	log = logger
}
