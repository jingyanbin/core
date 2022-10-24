package unicodeex

import (
	"github.com/jingyanbin/core/internal"
	"strconv"
)

// 转换为unicode编码
func ToRune(s string) rune {
	if len(s) < 3 || len(s) > 10 {
		panic(internal.NewError("ToRune error: %v", s))
	}
	if s[0] != '\\' || s[1] != 'u' {
		panic(internal.NewError("ToRune error: %v", s))
	}
	n, err := strconv.ParseUint(s[2:], 16, 32)
	if err != nil {
		panic(internal.NewError("ToRune error: %v, %v", s, err))
	}
	return rune(n)
}

// 字符编码范围
type RuneRange [2]rune

func (m *RuneRange) Check(r rune) bool {
	return r >= m[0] && r <= m[1]
}

type StringChecker []func(r rune) bool

func (m *StringChecker) runesCheck(s []rune) bool {
	for _, c := range s {
		ok := false
		for _, f := range *m {
			if f(c) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func (m *StringChecker) runesReplace(s []rune, rep rune) []rune {
	buf := make([]rune, 0, len(s))
	for _, c := range s {
		ok := false
		for _, f := range *m {
			if f(c) {
				ok = true
				break
			}
		}
		if !ok {
			buf = append(buf, rep)
		} else {
			buf = append(buf, c)
		}
	}
	return buf
}

func (m *StringChecker) Check(s string) bool {
	return m.runesCheck([]rune(s))
}

func (m *StringChecker) Replace(s string, rep rune) string {
	buf := m.runesReplace([]rune(s), rep)
	return string(buf)
}
