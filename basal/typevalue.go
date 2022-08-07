package basal

import "reflect"

type TypeValue struct {
	t reflect.Type
	v reflect.Value
}

func (m *TypeValue) SetFieldValueByName(value reflect.Value, name string) bool {
	var vElem reflect.Value
	switch m.v.Kind() {
	case reflect.Ptr:
		if m.v.IsNil() {
			return false
		}
		vElem = m.v.Elem()
	case reflect.Struct:
		vElem = m.v
	default:
		return false
	}
	field := vElem.FieldByName(name)
	if !field.CanSet() {
		return false
	}
	if field.Type() != value.Type() {
		return false
	}
	field.Set(value)
	return true
}

func (m *TypeValue) SetFieldValueByType(value reflect.Value) (ok bool) {
	m.RangeFields(func(tField reflect.StructField, vField reflect.Value) bool {
		if vField.CanSet() && vField.Type() == value.Type() {
			vField.Set(value)
			ok = true
			return false
		}
		return true
	})
	return
}

func (m *TypeValue) SetFieldByType(v interface{}) (ok bool) {
	return m.SetFieldValueByType(reflect.ValueOf(v))
}

func (m *TypeValue) SetFieldByName(v interface{}, name string) bool {
	return m.SetFieldValueByName(reflect.ValueOf(v), name)
}

func (m *TypeValue) RangeFields(f func(tField reflect.StructField, vField reflect.Value) bool) {
	var tElem reflect.Type
	var vElem reflect.Value
	switch m.v.Kind() {
	case reflect.Ptr:
		if m.v.IsNil() {
			return
		}
		tElem = m.t.Elem()
		vElem = m.v.Elem()
	case reflect.Struct:
		tElem = m.t
		vElem = m.v
	default:
		return
	}
	for i := 0; i < tElem.NumField(); i++ {
		tField := tElem.Field(i)
		vField := vElem.Field(i)
		if !f(tField, vField) {
			return
		}
	}
}

func (m *TypeValue) Type() reflect.Type {
	return m.t
}

func (m *TypeValue) Value() reflect.Value {
	return m.v
}

func (m *TypeValue) GetElem() (tElem reflect.Type, vElem reflect.Value, ok bool) {
	switch m.v.Kind() {
	case reflect.Ptr:
		if !m.v.IsNil() {
			tElem = m.t.Elem()
			vElem = m.v.Elem()
			ok = true
		}
	case reflect.Struct:
		tElem = m.t
		vElem = m.v
		ok = true
	}
	return
}

func (m *TypeValue) RangeMethods(f func(name string, method interface{}) bool) {
	for i := 0; i < m.v.NumMethod(); i++ {
		method := m.v.Method(i).Interface()
		name := m.t.Method(i).Name
		if !f(name, method) {
			return
		}
	}
}

func (m *TypeValue) GetMethodByName(name string) interface{} {
	method := m.v.MethodByName(name)
	if method.Kind() == reflect.Invalid {
		return nil
	}
	return method.Interface()
}

func TypeValueOf(ptr interface{}) *TypeValue {
	return &TypeValue{t: reflect.TypeOf(ptr), v: reflect.ValueOf(ptr)}
}
