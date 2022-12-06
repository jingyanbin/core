package log

import (
	"fmt"
	"github.com/jingyanbin/core/internal"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type LOGGER_LEVEL int

const (
	DEBUG LOGGER_LEVEL = 0
	INFO  LOGGER_LEVEL = 1
	WARN  LOGGER_LEVEL = 2
	ERROR LOGGER_LEVEL = 3
	FATAL LOGGER_LEVEL = 4
	OFF   LOGGER_LEVEL = 5
)

const logSkip = 2

func NewStdLogger(size int) *StdLogger {
	logger := &StdLogger{}
	logger.ch = make(chan *internal.Buffer, size)
	logger.wg.Add(1)
	go logger.run()
	return logger
}

var log = NewStdLogger(10000)

type StdLogger struct {
	level  LOGGER_LEVEL
	ch     chan *internal.Buffer
	closed int32
	wg     sync.WaitGroup
}

func (m *StdLogger) formatHeader(buf *internal.Buffer, level string, file string, line int, dt time.Time) {
	buf.AppendByte('[')
	buf.AppendString(level)
	buf.AppendByte(']')
	buf.AppendByte('[')
	buf.AppendInt(dt.Year(), 4)
	buf.AppendByte('-')
	buf.AppendInt(int(dt.Month()), 2)
	buf.AppendByte('-')
	buf.AppendInt(dt.Day(), 2)
	buf.AppendByte(' ')
	buf.AppendInt(dt.Hour(), 2)
	buf.AppendByte(':')
	buf.AppendInt(dt.Minute(), 2)
	buf.AppendByte(':')
	buf.AppendInt(dt.Second(), 2)
	buf.AppendByte(']')
	buf.AppendByte('[')
	buf.AppendString(file)
	buf.AppendBytes(':')
	buf.AppendInt(line, 0)
	buf.AppendBytes(']', ':')
}

func (m *StdLogger) run() {
	defer m.wg.Done()
	for v := range m.ch {
		m.write(v)
	}
}

func (m *StdLogger) write(buf *internal.Buffer) {
	defer internal.ExceptionError(nil)
	os.Stdout.Write(buf.Bytes())
	buf.Free()
}

func (m *StdLogger) push(buf *internal.Buffer) {
	if m.closed == 1 {
		return
	}
	defer internal.ExceptionError(nil)
	m.ch <- buf
}

func (m *StdLogger) output(level string, file string, line int, v ...interface{}) {
	var context string
	if len(v) > 1 {
		if format, ok := v[0].(string); ok {
			context = fmt.Sprintf(format, v[1:]...)
		} else {
			context = fmt.Sprint(v...)
		}
	} else {
		context = fmt.Sprint(v...)
	}
	buf := internal.NewBuffer(200 + len(context))
	m.formatHeader(buf, level, file, line, time.Now())
	buf.AppendString(context)
	buf.AppendByte('\n')
	m.push(buf)
}

func (m *StdLogger) Level() LOGGER_LEVEL {
	return m.level
}

func (m *StdLogger) Close() {
	if atomic.CompareAndSwapInt32(&m.closed, 0, 1) {
		close(m.ch)
	}
	m.wg.Wait()
}

func (m *StdLogger) SetLevel(level LOGGER_LEVEL) {
	if level < DEBUG || level > OFF {
		return
	}
	m.level = level
}

func (m *StdLogger) Debug(v ...any) {
	if m.level > DEBUG {
		return
	}
	file, line := internal.CallerShort(logSkip)
	m.output("D", file, line, v...)
}

func (m *StdLogger) Info(v ...any) {
	if m.level > INFO {
		return
	}
	file, line := internal.CallerShort(logSkip)
	m.output("I", file, line, v...)
}

func (m *StdLogger) Warn(v ...any) {
	if m.level > WARN {
		return
	}
	file, line := internal.CallerShort(logSkip)
	m.output("W", file, line, v...)
}

func (m *StdLogger) Error(v ...any) {
	if m.level > ERROR {
		return
	}
	file, line := internal.CallerShort(logSkip)
	m.output("E", file, line, v...)
}

func (m *StdLogger) Fatal(v ...any) {
	if m.level > FATAL {
		return
	}
	file, line := internal.CallerShort(logSkip)
	m.output("F", file, line, v...)
}
