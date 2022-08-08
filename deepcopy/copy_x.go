package deepcopy

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/vmihailenco/msgpack"
	"reflect"
)

func CopyByGob(src interface{}) interface{} {
	buf := &bytes.Buffer{}
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(src)
	if err != nil {
		panic(err)
	}
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	v := reflect.New(typ).Interface()
	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
	err = decoder.Decode(v)
	if err != nil {
		panic(err)
	}
	return v
}

func CopyByMsgPack(src interface{}) interface{} {
	data, err := msgpack.Marshal(src)
	if err != nil {
		panic(err)
	}
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	v := reflect.New(typ).Interface()
	err = msgpack.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
	return v
}

func CopyByJson(src interface{}) interface{} {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(src)
	if err != nil {
		panic(err)
	}
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	v := reflect.New(typ).Interface()
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	err = decoder.Decode(v)
	if err != nil {
		panic(err)
	}
	return v
}

// CopyByGoGo gogo protobuf 协议拷贝
func CopyByGoGo(src proto.Message) proto.Message {
	data, err := proto.Marshal(src)
	if err != nil {
		panic(err)
	}
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	v := reflect.New(typ).Interface().(proto.Message)
	err = proto.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
	return v
}
