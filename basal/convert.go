package basal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jingyanbin/core/internal"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

func IsUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				b := buf[i]
				for (b & 0x80) != 0 {
					b <<= 1  //左移一位
					nBytes++ //记录字符共占几个字节
				}
				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}
				nBytes-- //减掉首字节的一个计数
			}
		} else { //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

func Type(value interface{}) reflect.Type {
	return reflect.TypeOf(value)
}

func IsNilPointer(value interface{}) bool {
	vi := reflect.ValueOf(value)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func IsPointer(value interface{}) bool {
	vi := reflect.ValueOf(value)
	if vi.Kind() == reflect.Ptr {
		return true
	}
	return false
}

func SamePointer(pointers ...unsafe.Pointer) bool {
	if len(pointers) > 1 {
		p0 := *(*int)(pointers[0])
		for _, p := range pointers[1:] {
			if p0 != *(*int)(p) {
				return false
			}
		}
		return true
	}
	return false
}

func SamePtr(ptrs ...interface{}) bool {
	if len(ptrs) > 1 {
		vi0 := reflect.ValueOf(ptrs[0])
		if vi0.Kind() != reflect.Ptr {
			return false
		}
		for _, p := range ptrs[1:] {
			vix := reflect.ValueOf(p)
			if vix.Kind() != reflect.Ptr {
				return false
			}
			if vix.Pointer() != vi0.Pointer() {
				return false
			}
		}
		return true
	}
	return false
}

func ToJsonString(value interface{}, indent bool) (string, error) {
	b, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	if indent {
		var out bytes.Buffer
		err = json.Indent(&out, b, "", "    ")
		if err != nil {
			return string(b), nil
		}
		return out.String(), nil
	} else {
		return string(b), nil
	}
}

func ToString(value interface{}, indent bool) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int, int8, int16, int32, int64:
		n, err := ToInt64(v)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(n, 10), nil
	case uint, uint8, uint16, uint32, uint64:
		n, err := ToInt64(v)
		if err != nil {
			return "", err
		}
		return strconv.FormatUint(uint64(n), 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case map[string]interface{}, []interface{}:
		return ToJsonString(v, indent)
	case []byte:
		return string(v), nil
	case json.Number:
		return string(v), nil
	default:
		vi := reflect.ValueOf(value)
		kd := vi.Kind()
		switch kd {
		case reflect.Struct:
			return ToJsonString(vi, indent)
		case reflect.Ptr:
			if vi.IsNil() {
				return fmt.Sprintf("<nil %v>", vi.Type()), nil
			}
			kd2 := vi.Elem().Kind()
			switch kd2 {
			case reflect.Struct:
				return ToJsonString(vi.Elem().Interface(), indent)
			default:
				return "", NewError("ToString value ptr type error: %v", vi.Type())
			}
		default:
			return "", NewError("ToString value type error: %v", vi.Type())
		}
	}
}

func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case []byte:
		return strconv.ParseFloat(string(v), 64)
	case json.Number:
		return v.Float64()
	default:
		switch value := reflect.ValueOf(v); value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n := value.Int()
			return float64(n), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n := value.Uint()
			return float64(n), nil
		case reflect.Float64, reflect.Float32:
			return value.Float(), nil
		case reflect.String:
			return strconv.ParseFloat(value.String(), 64)
		case reflect.Slice:
			return strconv.ParseFloat(string(value.Bytes()), 64)
		}
	}
	return 0, NewError("ToFloat64 value type error: %v", Type(value))
}

func ToFloat32(value interface{}) (float32, error) {
	v, err := ToFloat64(value)
	return float32(v), err
}

func ToInt64(value interface{}) (int64, error) {
	switch n := value.(type) {
	case bool:
		if n {
			return 1, nil
		} else {
			return 0, nil
		}
	case int:
		return int64(n), nil
	case int8:
		return int64(n), nil
	case int16:
		return int64(n), nil
	case int32:
		return int64(n), nil
	case int64:
		return n, nil
	case uint:
		return int64(n), nil
	case uint8:
		return int64(n), nil
	case uint16:
		return int64(n), nil
	case uint32:
		return int64(n), nil
	case uint64:
		return int64(n), nil
	case float64:
		return int64(n), nil
	case float32:
		return int64(n), nil
	case string:
		f, err := strconv.ParseFloat(n, 64)
		return int64(f), err
		//return strconv.ParseInt(n, 10, 64)
	case []byte:
		f, err := strconv.ParseFloat(string(n), 64)
		return int64(f), err
		//return strconv.ParseInt(string(n), 10, 64)
	case json.Number:
		return n.Int64()
	default:
		switch value := reflect.ValueOf(n); value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return value.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return int64(value.Uint()), nil
		case reflect.Float64, reflect.Float32:
			return int64(value.Float()), nil
		case reflect.String:
			f, err := strconv.ParseFloat(value.String(), 64)
			return int64(f), err
			//return strconv.ParseInt(value.String(), 10, 64)
		case reflect.Slice:
			f, err := strconv.ParseFloat(string(value.Bytes()), 64)
			return int64(f), err
			//return strconv.ParseInt(string(value.Bytes()), 10, 64)
		}
	}
	return 0, NewError("ToInt64 value type error: %v", Type(value))
}

func ToInt32(value interface{}) (int32, error) {
	v, err := ToInt64(value)
	return int32(v), err
}

func ToInt16(value interface{}) (int16, error) {
	v, err := ToInt64(value)
	return int16(v), err
}

func ToInt8(value interface{}) (int8, error) {
	v, err := ToInt64(value)
	return int8(v), err
}

func ToInt(value interface{}) (int, error) {
	v, err := ToInt64(value)
	return int(v), err
}

func ToUint64(value interface{}) (uint64, error) {
	v, err := ToInt64(value)
	return uint64(v), err
}

func ToUint32(value interface{}) (uint32, error) {
	v, err := ToInt64(value)
	return uint32(v), err
}

func ToUint16(value interface{}) (uint16, error) {
	v, err := ToInt64(value)
	return uint16(v), err
}

func ToUint8(value interface{}) (uint8, error) {
	v, err := ToInt64(value)
	return uint8(v), err
}

func ToUint(value interface{}) (uint, error) {
	v, err := ToInt64(value)
	return uint(v), err
}

func Int64ToBool(value int64) bool {
	return value != 0
}

func Int32ToBool(value int32) bool {
	return value != 0
}

// 驼峰写法转下划线小写 eg: LevelAbc=>level_abc
func ToLowerLine(s string) string {
	s2 := make([]byte, 0, len(s)+1)
	for i, c := range []byte(s) {
		if c > 64 && c < 91 {
			if i > 0 {
				s2 = append(s2, '_')
			}
			s2 = append(s2, c+32)
		} else {
			s2 = append(s2, c)
		}
	}
	return string(s2)
}

// 字符串转bytes 慎修改转换后的值
func StrPtrToBytes(s string) []byte {
	p := *(*[2]uintptr)(unsafe.Pointer(&s))
	p2 := [3]uintptr{p[0], p[1], p[1]}
	return *(*[]byte)(unsafe.Pointer(&p2))
}

// bytes转字符串
func BytesPtrToStr(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// 字符串中字符前(不存在此字符就加此字符) 一般mysql语句拼接使用 eg: old '\”, before add '\\'
func StrAddBeforeNotHas(s string, old rune, add rune) string {
	data := []rune(s)
	bs := make([]rune, 0, len(data))
	for _, r := range data {
		if r == old {
			bs = append(bs, add, old)
		} else {
			bs = append(bs, r)
		}
	}
	return string(bs)
}

//const UINT64_MIN uint64 = 0
//const UINT64_MAX uint64 = ^UINT64_MIN
//const INT64_MIN = ^UINT64_MAX
//const INT64_MAX  = int64(^uint64(0)>>1)

const UINT32_MIN uint32 = 0
const UINT32_MAX = ^UINT32_MIN
const INT32_MIN = ^UINT32_MAX
const INT32_MAX = int32(^uint32(0) >> 1)

const UINT16_MIN uint16 = 0
const UINT16_MAX = ^UINT16_MIN
const INT16_MIN = ^UINT16_MAX
const INT16_MAX = int16(UINT16_MAX >> 1)

const overfolw63div10 = (1<<63 - 1) / 10

func AtoInt64(s string) (x int64, err error) {
	neg := false
	if s == "" {
		return 0, NewError("param error: %s", s)
	}

	if s[0] == '-' || s[0] == '+' {
		neg = s[0] == '-'
		s = s[1:]
	} else if s[0] < '0' || s[0] > '9' {
		return 0, NewError("param error: %s", s)
	}

	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > overfolw63div10 {
			// overflow
			return 0, NewError("param error: overflow %v", s)
		}
		x = x*10 + int64(c) - '0'
		if x < 0 {
			// overflow
			return 0, NewError("param error: overflow %v", s)
		}
	}
	if neg {
		x = -x
	}
	return x, nil
}

func AtoInt(s string) (int, error) {
	x, err := AtoInt64(s)
	return int(x), err
}

func AbsInt64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

func AbsInt32(n int32) int32 {
	y := n >> 31
	return (n ^ y) - y
}

func AbsInt16(n int16) int16 {
	y := n >> 15
	return (n ^ y) - y
}

func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func AbsFloat32(n float32) float32 {
	if n < 0 {
		return -n
	}
	return n
}

func AbsFloat64(n float64) float64 {
	if n < 0 {
		return -n
	}
	return n
}

func Round(value float64, digit int) float64 {
	p10 := math.Pow10(digit)
	return math.Trunc((value+0.5/p10)*p10) / p10
}

func AddRemainInt64(oldNum, addNum, numMax int64) (newNum int64, added int64, remained int64) {
	if addNum < 0 {
		return oldNum, 0, 0
	}
	cha := numMax - oldNum
	remained = addNum - cha
	if remained > 0 {
		return numMax, addNum - remained, remained
	}
	return oldNum + addNum, addNum, 0
}

func AddRemainInt32(oldNum, addNum, numMax int32) (newNum int32, added int32, remained int32) {
	if addNum < 0 {
		return oldNum, 0, 0
	}
	cha := numMax - oldNum
	remained = addNum - cha
	if remained > 0 {
		return numMax, addNum - remained, remained
	}
	return oldNum + addNum, addNum, 0
}

type OnceSuccess = internal.OnceSuccess

type Waiter = internal.Waiter
