package basal

import (
	"reflect"
	"runtime"
	"strings"
)

func GetFuncName(i interface{}, seps ...rune) string {
	// 获取函数名称
	fn := GetFuncFullName(i)
	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})
	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}

func GetFuncFullName(i interface{}) string {
	if i == nil {
		return "nil"
	}
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func GetFuncShortName(i interface{}) string {
	name := GetFuncFullName(i)
	if len(name) > 0 {
		names := strings.Split(name, ".")
		dLen := len(names)
		if dLen > 0 {
			return names[len(names)-1]
		}
	}
	return name
}
