package basal

import (
	"strconv"
)

// 转换为unicode编码
func ToRune(s string) rune {
	if len(s) < 3 || len(s) > 10 {

		panic(NewError("ToRune error: %v", s))
	}
	if s[0] != '\\' || s[1] != 'u' {
		panic(NewError("ToRune error: %v", s))
	}
	n, err := strconv.ParseUint(s[2:], 16, 32)
	if err != nil {
		panic(NewError("ToRune error: %v, %v", s, err))
	}
	return rune(n)
}

// 字符编码范围
type RuneRange [2]rune

func (m *RuneRange) Check(r rune) bool {
	return r >= m[0] && r <= m[1]
}

// 字符串编码检查
type UnicodeChecker []func(r rune) bool

func (m *UnicodeChecker) runesCheck(s []rune) bool {
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

func (m *UnicodeChecker) runesReplace(s []rune, rep rune) []rune {
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

func (m *UnicodeChecker) Check(s string) bool {
	return m.runesCheck([]rune(s))
}

func (m *UnicodeChecker) Replace(s string, rep rune) string {
	buf := m.runesReplace([]rune(s), rep)
	return string(buf)
}

type LenChecker func(r rune) int

func (m *LenChecker) Len(s string) (n int) {
	for _, c := range s {
		n += (*m)(c)
	}
	return
}

func (m *LenChecker) Check(s string, max int) bool {
	if n := m.Len(s); n > max {
		return false
	}
	return true
}
