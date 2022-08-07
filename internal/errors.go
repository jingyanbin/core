package internal

import (
	"errors"
	"fmt"
	_ "unsafe"
)

//go:linkname Sprintf github.com/jingyanbin/core/basal.Sprintf
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

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
