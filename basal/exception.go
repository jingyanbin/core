package basal

import "fmt"

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func NewError(format string, a ...interface{}) error

func ToError(r interface{}) (err error)

func CallerShort(skip int) (file string, line int)

func Caller(skip int) (file string, line int)

func CallerInFunc(skip int) (name string, file string, line int)

func CallerLineStack(stack string) (name string, file string)

func ExceptionError(catch func(e error))

func Exception(catch ...func(stack string, e error))

func Try(f func(), catch func(stack string, e error))
