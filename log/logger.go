package log

import (
	"fmt"
	"github.com/jingyanbin/core/datetime"
	internal "github.com/jingyanbin/core/internal"
	tz "github.com/jingyanbin/core/timezone"
	"io"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	_ "unsafe"
)

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
	LOG_LEVEL_FATAL
	LOG_LEVEL_OFF
)

const logSkip = 3
const logChanSize = 1024

var logLevels = []string{"DEBU", "INFO", "WARN", "ERRO", "CRIT"}

type logChan chan *logItem

func (ch logChan) Push(item *logItem) (succeed bool) {
	defer internal.Exception(nil)
	ch <- item
	return true
}

type logWriter interface {
	Write(unix int64, level int, file string, line int, content string)
	SetFormatHeader(formatHeader func(buf *internal.Buffer, level string, line int, file string, dt *datetime.DateTime))
	Close()
}

type logBaseWriter struct {
	time         datetime.DateTime
	customHeader func(buf *internal.Buffer, level string, line int, file string, dt *datetime.DateTime)
}

func (my *logBaseWriter) SetZone(zone tz.TimeZone) {
	my.time.SetZone(zone)
}

func (my *logBaseWriter) Close() {}

func (my *logBaseWriter) SetFormatHeader(formatHeader func(buf *internal.Buffer, level string, line int, file string, dt *datetime.DateTime)) {
	my.customHeader = formatHeader
}

func (my *logBaseWriter) formatHeader(buf *internal.Buffer, level string, line int, file string) {
	buf.AppendBytes('[')
	buf.AppendString(level)
	buf.AppendBytes(' ')
	buf.AppendInt(my.time.Year(), 4)
	buf.AppendBytes('/')
	buf.AppendInt(my.time.Month(), 2)
	buf.AppendBytes('/')
	buf.AppendInt(my.time.Day(), 2)
	buf.AppendBytes(' ')
	buf.AppendInt(my.time.Hour(), 2)
	buf.AppendBytes(':')
	buf.AppendInt(my.time.Min(), 2)
	buf.AppendBytes(':')
	buf.AppendInt(my.time.Sec(), 2)
	buf.AppendBytes(' ')
	buf.AppendString(file)
	buf.AppendBytes(':')
	buf.AppendInt(line, 0)
	buf.AppendString("]")
}

func (my *logBaseWriter) write(writer io.Writer, level int, file string, line int, content string) {
	//buf := make([]byte, 0, 40+len(file)+len(content))
	buf := internal.NewBuffer(40 + len(file) + len(content))
	defer buf.Free()
	if my.customHeader == nil {
		my.formatHeader(buf, logLevels[level], line, file)
	} else {
		my.customHeader(buf, logLevels[level], line, file, &my.time)
	}
	buf.AppendByte(' ')
	buf.AppendString(content)
	buf.AppendByte('\n')
	//buf = append(buf, content...)
	//buf = append(buf, '\n')
	n, err := writer.Write(buf.Bytes())

	if err != nil {
		fmt.Printf("logBaseWriter write error: %v, n=%v\n", err, n)
	}
}

type logStdWriter struct {
	logBaseWriter
	writer io.Writer
}

func (my *logStdWriter) Write(unix int64, level int, file string, line int, content string) {
	my.time.FlushToUnix(unix)
	my.write(my.writer, level, file, line, content)
}

func NewLogStdWriter(writer io.Writer) *logStdWriter {
	out := &logStdWriter{}
	out.writer = writer
	out.time.SetZone(tz.Local())
	return out
}

type logFileWriter struct {
	logBaseWriter
	writer            *internal.HandleFile
	filePathFormatter string
	mu                sync.Mutex
}

func (my *logFileWriter) NextPathName() (folderPath string, fileName string) {
	return filepath.Split(my.time.Format(my.filePathFormatter))
}

func (my *logFileWriter) Write(unix int64, level int, file string, line int, content string) {
	my.mu.Lock()
	defer my.mu.Unlock()
	my.time.FlushToUnix(unix)
	my.writer.SetPathName(my.NextPathName())
	my.write(my.writer, level, file, line, content)
}

func (my *logFileWriter) Close() {
	my.mu.Lock()
	defer my.mu.Unlock()
	my.writer.Close()
}

func NewLogFileWriter(filePathFormatter string) *logFileWriter {
	out := &logFileWriter{}
	out.writer = internal.NewHandleFile(internal.HANDLE_FILE_FLAG_WRITER, internal.HANDLE_FILE_PERM_ALL)
	out.time.SetZone(tz.Local())
	if filePathFormatter == "" {
		out.filePathFormatter = internal.Path.ProgramDirJoin("output.log.%Y-%m-%d-%H")
	} else {
		out.filePathFormatter = filePathFormatter
	}
	return out
}

type stdLogger struct {
	level    int
	handlers []logWriter
	q        logChan
	running  int32
	async    bool
	wg       sync.WaitGroup
}

func (my *stdLogger) output(level int, v ...interface{}) {
	if level < my.level {
		return
	}
	file, line := internal.CallerShort(logSkip)
	my.write(internal.Unix(), level, file, line, fmt.Sprint(v...))
}

func (my *stdLogger) outputf(level int, format string, v ...interface{}) {
	if level < my.level {
		return
	}
	file, line := internal.CallerShort(logSkip)
	my.write(internal.Unix(), level, file, line, fmt.Sprintf(format, v...))
}

//func (my *stdLogger) output(level int, content string) {
//	if level < my.level {
//		return
//	}
//	file, line := CallerShort(logSkip)
//	my.write(Unix(), level, file, line, content)
//}

func (my *stdLogger) SetLevel(level int) {
	if level < LOG_LEVEL_DEBUG || level > LOG_LEVEL_OFF {
		return
	}
	my.level = level
}

func (my *stdLogger) AddHandler(handlers ...logWriter) {
	for _, handler := range handlers {
		if handler == nil {
			continue
		}
		my.handlers = append(my.handlers, handler)
	}
}

func (my *stdLogger) SetAsync(async bool) {
	my.async = async
	if async {
		my.start()
	} else {
		my.Wait()
	}
}

func (my *stdLogger) start() {
	if atomic.CompareAndSwapInt32(&my.running, 0, 1) {
		my.q = make(logChan, logChanSize)
		my.wg.Add(1)
		go my.run()
	}
}

func (my *stdLogger) Wait() {
	if atomic.CompareAndSwapInt32(&my.running, 1, 0) {
		close(my.q)
		my.wg.Wait()
	}
	my.close()
}

func (my *stdLogger) SetFormatHeader(formatHeader func(buf *internal.Buffer, level string, line int, file string, dt *datetime.DateTime)) {
	for _, handler := range my.handlers {
		handler.SetFormatHeader(formatHeader)
	}
}

func (my *stdLogger) close() {
	for _, handler := range my.handlers {
		handler.Close()
	}
}

func (my *stdLogger) run() {
	//defer Exception(func(stack string, e error) {
	//	panic(NewError("stdLogger run panic: %s \n%s", e, stack))
	//})
	defer my.wg.Done()
	var handler logWriter
	for item := range my.q {
		for _, handler = range my.handlers {
			handler.Write(item.unix, item.level, item.file, item.line, item.content)
		}
		item.free()
	}
}

func (my *stdLogger) write(unix int64, level int, file string, line int, content string) {
	if my.async {
		item := logItemFree.Get().(*logItem)
		item.unix = unix
		item.level = level
		item.content = content
		item.file = file
		item.line = line
		if my.q.Push(item) {
			return
		}
	}
	for _, handler := range my.handlers {
		handler.Write(unix, level, file, line, content)
	}
}

func (my *stdLogger) Debug(v ...interface{}) {
	my.output(LOG_LEVEL_DEBUG, v...)
}
func (my *stdLogger) Info(v ...interface{}) {
	my.output(LOG_LEVEL_INFO, v...)
}
func (my *stdLogger) Warn(v ...interface{}) {
	my.output(LOG_LEVEL_WARN, v...)
}
func (my *stdLogger) Error(v ...interface{}) {
	my.output(LOG_LEVEL_ERROR, v...)
}
func (my *stdLogger) Fatal(v ...interface{}) {
	my.output(LOG_LEVEL_FATAL, v...)
}

func (my *stdLogger) DebugF(format string, v ...interface{}) {
	my.outputf(LOG_LEVEL_DEBUG, format, v...)
}
func (my *stdLogger) InfoF(format string, v ...interface{}) {
	my.outputf(LOG_LEVEL_INFO, format, v...)
}
func (my *stdLogger) WarnF(format string, v ...interface{}) {
	my.outputf(LOG_LEVEL_WARN, format, v...)
}
func (my *stdLogger) ErrorF(format string, v ...interface{}) {
	my.outputf(LOG_LEVEL_ERROR, format, v...)
}
func (my *stdLogger) FatalF(format string, v ...interface{}) {
	my.outputf(LOG_LEVEL_FATAL, format, v...)
}

func NewLogger(level int, async bool, handlers ...logWriter) *stdLogger {
	logger := &stdLogger{}
	logger.SetLevel(level)
	logger.AddHandler(handlers...)
	logger.SetAsync(async)
	//runtime.SetFinalizer(logger, (*stdLogger).Wait)
	loggerMgr.Append(logger)
	return logger
}

type loggerManager struct {
	loggers []*stdLogger
	mu      sync.Mutex
}

func (my *loggerManager) Append(logger *stdLogger) {
	my.mu.Lock()
	defer my.mu.Unlock()
	my.loggers = append(my.loggers, logger)
}

func (my *loggerManager) Wait() {
	for _, logger := range my.loggers {
		logger.Wait()
	}
}

type logItem struct {
	unix    int64
	level   int
	content string
	file    string
	line    int
}

func (it *logItem) free() {
	logItemFree.Put(it)
}

var loggerMgr loggerManager

var logItemFree = &sync.Pool{New: func() interface{} { return new(logItem) }}

var log = &stdLogger{level: LOG_LEVEL_DEBUG, handlers: []logWriter{NewLogStdWriter(os.Stdout)}}

func init() {
	log.SetAsync(true)
}
