package internal

import (
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	_ "unsafe"
)

//go:linkname NewError github.com/jingyanbin/core/basal.NewError
func NewError(format string, a ...interface{}) error {
	return errors.New(fmt.Sprintf(format, a...))
}

//go:linkname ToError github.com/jingyanbin/core/basal.ToError
func ToError(r interface{}) (err error) {
	switch x := r.(type) {
	case string:
		err = NewError(x)
	case error:
		err = x
	default:
		err = NewError("unknown error: %v", x)
	}
	return
}

const exceptionSkip = 3

// 调用信息短文件名
//
//go:linkname CallerShort github.com/jingyanbin/core/basal.CallerShort
func CallerShort(skip int) (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(skip)
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	} else {
		file = "???"
		line = 0
	}
	return
}

// 调用信息长文件名
//
//go:linkname Caller github.com/jingyanbin/core/basal.Caller
func Caller(skip int) (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}
	return
}

//go:linkname CallerInFunc github.com/jingyanbin/core/basal.CallerInFunc
func CallerInFunc(skip int) (name string, file string, line int) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(skip)
	if ok {
		inFunc := runtime.FuncForPC(pc)
		name = inFunc.Name()
	} else {
		file = "???"
		name = "???"
	}
	return
}

//go:linkname CallerLineStack github.com/jingyanbin/core/basal.CallerLineStack
func CallerLineStack(stack string) (name string, file string) {
	name = "???"
	file = "???"
	stackArr := strings.SplitN(stack, "panic.go", 2)
	if len(stackArr) != 2 {
		return
	}
	stackLines := strings.SplitN(stackArr[1], "\n", 4)
	if len(stackLines) != 4 {
		return
	}
	//name = strings.Trim(stackLines[1], "\n")
	//file = strings.Trim(stackLines[2], "\n")
	name = strings.TrimSpace(stackLines[1])
	file = strings.TrimSpace(stackLines[2])
	return
}

func formatStack(name, file string, err string, stack []byte) *Buffer {
	buf := NewBuffer(160 + len(stack) + len(name))
	buf.AppendStrings("exception: ", err, "\nfile: ", file)
	buf.AppendStrings("\nfunc: ", name, "\n")
	buf.AppendBytes(stack...)
	return buf
}

//go:linkname Exception github.com/jingyanbin/core/basal.Exception
func Exception(catch ...func(stack string, e error)) {
	if err := recover(); err != nil {
		if len(catch) == 0 {
			return
		}
		info := debug.Stack()
		name, file := CallerLineStack(string(info))
		myErr := ToError(err)
		myBuf := formatStack(name, file, myErr.Error(), info)
		defer myBuf.Free()
		for _, f := range catch {
			if f == nil {
				continue
			}
			f(myBuf.ToString(), myErr)
		}
	}
}

//go:linkname ExceptionError github.com/jingyanbin/core/basal.ExceptionError
func ExceptionError(catch func(e error)) {
	if err := recover(); err != nil {
		catch(ToError(err))
	}
}

//go:linkname Try github.com/jingyanbin/core/basal.Try
func Try(f func(), catch func(stack string, e error)) {
	defer Exception(catch)
	f()
}
