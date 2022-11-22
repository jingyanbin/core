package internal

import (
	"bytes"
	//json "encoding/json"
	json "github.com/json-iterator/go"
	"os"
	_ "unsafe"
)

type JsonINumber interface {
	String() string
	Float64() (float64, error)
	Int64() (int64, error)
}

//go:linkname LoadJsonFileTo github.com/jingyanbin/core/jsonex.LoadFileTo
func LoadJsonFileTo(jsFileName string, toPtr interface{}) error {
	//if isNil, typ := IsNilPointer(v); isNil {
	//	return NewError("LoadJsonFileTo IsNilPointer: %s", typ)
	//}
	data, err := os.ReadFile(jsFileName)
	if err != nil {
		return err
	}
	return LoadJsonBytesTo(data, toPtr)
}

//go:linkname LoadJsonBytesTo github.com/jingyanbin/core/jsonex.LoadBytesTo
func LoadJsonBytesTo(js []byte, toPtr interface{}) error {
	decoder := json.NewDecoder(bytes.NewBuffer(js))
	//decoder := json.NewDecoder(bytes.NewBuffer(js))
	decoder.UseNumber()
	err := decoder.Decode(toPtr)
	if err != nil {
		return err
	}
	return nil
}

//go:linkname LoadJsonStringTo github.com/jingyanbin/core/jsonex.LoadStringTo
func LoadJsonStringTo(js string, toPtr interface{}) error {
	return LoadJsonBytesTo([]byte(js), toPtr)
}

//go:linkname TryDumpJson github.com/jingyanbin/core/jsonex.TryDump
func TryDumpJson(v interface{}, indent bool) (string, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	if indent {
		enc.SetIndent("", "\t")
	}
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	return buf.String(), err
}

//go:linkname DumpJson github.com/jingyanbin/core/jsonex.Dump
func DumpJson(v interface{}, indent bool) string {
	s, err := TryDumpJson(v, indent)
	if err != nil {
		panic(err)
	}
	return s
}
