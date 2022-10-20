package basal

type StringChecker []func(r rune) bool

func (m *StringChecker) runesCheck(s []rune) (int, bool) {
	for i, c := range s {
		success := false
		for _, f := range *m {
			if !success && f(c) {
				success = true
				break
			}
		}
		if !success {
			return i, false
		}
	}
	return 0, true
}

func (m *StringChecker) Check(s string) bool {
	_, ok := m.runesCheck([]rune(s))
	return ok
}

func (m *StringChecker) CheckIndex(s string) (int, bool) {
	str := []rune(s)
	if index, ok := m.runesCheck(str); ok {
		return index, true
	} else {
		return index, false
	}
}
