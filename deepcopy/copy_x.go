package deepcopy

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"github.com/vmihailenco/msgpack"
	"reflect"
)

// CopyByGob[T any]
//
//	@Description: 不能拷贝私有字段
//	@param src
//	@return T
func CopyByGob[T any](src T) T {
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
	v := reflect.New(typ).Interface().(T)
	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
	err = decoder.Decode(v)
	if err != nil {
		panic(err)
	}
	return v
}

// CopyByMsgPack[T any]
//
//	@Description: 不能拷贝私有字段
//	@param src
//	@return T
func CopyByMsgPack[T any](src T) T {
	data, err := msgpack.Marshal(src)
	if err != nil {
		panic(err)
	}
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	v := reflect.New(typ).Interface().(T)
	err = msgpack.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
	return v
}

// CopyByJson[T any]
//
//	@Description: 不能拷贝私有字段
//	@param src
//	@return T
func CopyByJson[T any](src T) T {
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
	v := reflect.New(typ).Interface().(T)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	err = decoder.Decode(v)
	if err != nil {
		panic(err)
	}
	return v
}

// CopyByGoGo[T proto.Message]
//
//	@Description: 协议拷贝 不能拷贝私有字段
//	@param src
//	@return T
func CopyByGoGo[T proto.Message](src T) T {
	data, err := proto.Marshal(src)
	if err != nil {
		panic(err)
	}
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	v := reflect.New(typ).Interface().(T)
	err = proto.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
	return v
}
