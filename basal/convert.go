package basal

import (
	"github.com/jingyanbin/core/internal"
	"reflect"
	"unsafe"
)

func IsUTF8(buf []byte) bool {
	return internal.IsUTF8(buf)
}

func Type(value interface{}) reflect.Type {
	return internal.Type(value)
}

func IsNilPointer(value interface{}) bool {
	return internal.IsNilPointer(value)
}

func IsPointer(value interface{}) bool {
	return internal.IsPointer(value)
}

func SamePointer(pointers ...unsafe.Pointer) bool {
	return internal.SamePointer(pointers...)
}

func SamePtr(ptrs ...interface{}) bool {
	return internal.SamePtr(ptrs...)
}

func ToString(value interface{}, indent bool) (string, error) {
	return internal.ToString(value, indent)
}

func ToFloat64(value interface{}) (float64, error) {
	return internal.ToFloat64(value)
}

func ToFloat32(value interface{}) (float32, error) {
	return internal.ToFloat32(value)
}

func ToInt64(value interface{}) (int64, error) {
	return internal.ToInt64(value)
}

func ToInt32(value interface{}) (int32, error) {
	return internal.ToInt32(value)
}

func ToInt16(value interface{}) (int16, error) {
	return internal.ToInt16(value)
}

func ToInt8(value interface{}) (int8, error) {
	return internal.ToInt8(value)
}

func ToInt(value interface{}) (int, error) {
	return internal.ToInt(value)
}

func ToUint64(value interface{}) (uint64, error) {
	return internal.ToUint64(value)
}

func ToUint32(value interface{}) (uint32, error) {
	return internal.ToUint32(value)
}

func ToUint16(value interface{}) (uint16, error) {
	return internal.ToUint16(value)
}

func ToUint8(value interface{}) (uint8, error) {
	return internal.ToUint8(value)
}

func ToUint(value interface{}) (uint, error) {
	return internal.ToUint(value)
}

func NumberToBool[T Number](value T) bool {
	return internal.NumberToBool(value)
}

// 驼峰写法转下划线小写 eg: LevelAbc=>level_abc
func ToLowerLine(s string) string {
	return internal.ToLowerLine(s)
}

// 字符串转bytes 慎修改转换后的值
func StrPtrToBytes(s string) []byte {
	return internal.StrPtrToBytes(s)
}

// bytes转字符串
func BytesPtrToStr(bs []byte) string {
	return internal.BytesPtrToStr(bs)
}

// 字符串中字符前(不存在此字符就加此字符) 一般mysql语句拼接使用 eg: old '\”, before add '\\'
func StrAddBeforeNotHas(s string, old rune, add rune) string {
	return internal.StrAddBeforeNotHas(s, old, add)
}

// 有小数的直接忽略
func AtoInt64(s string) (x int64, err error) {
	return internal.AtoInt64(s)
}

func AtoInt32(s string) (int32, error) {
	return internal.AtoInt32(s)
}

func AtoInt(s string) (int, error) {
	return internal.AtoInt(s)
}

func AbsInt64(n int64) int64 {
	return internal.AbsInt64(n)
}

func AbsInt32(n int32) int32 {
	return internal.AbsInt32(n)
}

func AbsInt16(n int16) int16 {
	return internal.AbsInt16(n)
}

func AbsInt8(n int8) int8 {
	return internal.AbsInt8(n)
}

func Abs[T Number](n T) T {
	return internal.Abs(n)
}

func Round(value float64, digit int) float64 {
	return internal.Round(value, digit)
}

// 添加剩余
func AddRemain[T Number](oldNum, addNum, numMax T) (newNum, added, remained T) {
	return internal.AddRemain(oldNum, added, numMax)
}
