package internal

import (
	"fmt"
	"os"
	_ "unsafe"
)

type ILogger interface {
	InfoF(format string, v ...interface{})
	ErrorF(format string, v ...interface{})
}

type stdoutLogger struct {
	ch chan string
}

func (m *stdoutLogger) run() {
	for {
		select {
		case s, ok := <-m.ch:
			if !ok {
				panic("basal stdoutLogger error")
			}
			n, err := os.Stdout.WriteString(s)
			if n == 0 {
				panic(Sprintf("basal stdoutLogger WriteString n is 0"))
			} else if err != nil {
				panic(Sprintf("basal stdoutLogger err: %v", err))
			}
		}
	}
}

func (m *stdoutLogger) output(level string, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	sLen := len(s)
	if (sLen > 0 && s[sLen-1] != '\n') || sLen == 0 {
		m.ch <- level + s + "\n"
	} else {
		m.ch <- level + s
	}
}

func (m *stdoutLogger) ErrorF(format string, v ...interface{}) {
	m.output("[ERRO] ", format, v...)
}

func (m *stdoutLogger) InfoF(format string, v ...interface{}) {
	m.output("[INFO] ", format, v...)
}

func newStdoutLogger() *stdoutLogger {
	logger := &stdoutLogger{ch: make(chan string, 10000)}
	go logger.run()
	return logger
}

var stdLogger ILogger = newStdoutLogger()
var log ILogger = stdLogger

//go:linkname GetStdoutLogger github.com/jingyanbin/core/basal.GetStdoutLogger
func GetStdoutLogger() ILogger {
	return stdLogger
}

func SetLogger(logger ILogger) {
	log = logger
}
