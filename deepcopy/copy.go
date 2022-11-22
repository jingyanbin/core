package deepcopy

import (
	"fmt"
	"github.com/jingyanbin/core/internal"
	"reflect"
	"unsafe"
)

var log = internal.GetStdoutLogger()

type flag uintptr

const (
	flagKindWidth        = 5 // there are 27 kinds
	flagKindMask    flag = 1<<flagKindWidth - 1
	flagStickyRO    flag = 1 << 5
	flagEmbedRO     flag = 1 << 6
	flagIndir       flag = 1 << 7
	flagAddr        flag = 1 << 8
	flagMethod      flag = 1 << 9
	flagMethodShift      = 10
	flagRO          flag = flagStickyRO | flagEmbedRO
)

type myValue struct {
	typ *struct{}
	ptr unsafe.Pointer
	flag
}

func checkCanSet(dst reflect.Value, all bool) bool {
	if all {
		return true
	}
	return dst.CanSet()
}

func setValue(dst, src reflect.Value, all bool) bool {
	if dst.CanSet() {
		dst.Set(src)
		return true
	}
	if all {
		dstP := *(*myValue)(unsafe.Pointer(&dst))
		dstP.flag = dstP.flag ^ (flagAddr | flagRO)
		dstV := *(*reflect.Value)(unsafe.Pointer(&dstP))
		dstV.Set(src)
		return true
	}
	return false
}

// Copy 拷贝结构体公共字段, 不支持chan
func Copy(src interface{}) interface{} {
	return deepCopy(src, false)
}

// CopyAll 拷贝结构体所有字段(包含私有字段), 不支持chan
func CopyAll(src interface{}) interface{} {
	return deepCopy(src, true)
}

func deepCopy(src interface{}, all bool) interface{} {
	if src == nil {
		return nil
	}
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.New(srcValue.Type()).Elem()
	parent := make(map[uintptr]struct{})
	copyValue(dstValue, srcValue, all, parent)
	return dstValue.Interface()
}

func copyValue(dst, src reflect.Value, all bool, parent map[uintptr]struct{}) {
	//if src.CanInterface() {
	//	if !checkCanSet(dst, all) {
	//		return
	//	}
	//	if all {
	//		if copier, ok := src.Interface().(IDeepCopyAll); ok {
	//			setValue(dst, reflect.ValueOf(copier.DeepCopyAll()), all)
	//			return
	//		}
	//	} else {
	//		if copier, ok := src.Interface().(IDeepCopy); ok {
	//			setValue(dst, reflect.ValueOf(copier.DeepCopy()), all)
	//			return
	//		}
	//	}
	//}
	switch src.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		if !checkCanSet(dst, all) {
			return
		}
		if _, ok := parent[src.Pointer()]; ok {
			panic(fmt.Sprintf("指针循环引用: reflect.Ptr: %v\n", src))
			return
		}
		parent[src.Pointer()] = struct{}{}
		srcValue := src.Elem()
		dstNewPtr := reflect.New(srcValue.Type())
		copyValue(dstNewPtr.Elem(), srcValue, all, parent)
		setValue(dst, dstNewPtr, all)

		delete(parent, src.Pointer())
		//parent[src.Pointer()] = struct{}{}

		//dstP := *(*myValue)(unsafe.Pointer(&dst))
		//srcP := *(*myValue)(unsafe.Pointer(&dstValuePtr))
		//*(*unsafe.Pointer)(dstP.ptr) = srcP.ptr
	case reflect.Interface:
		if src.IsNil() {
			return
		}
		if !checkCanSet(dst, all) {
			return
		}
		//dstValue := dst.Elem()
		srcValue := src.Elem()
		dstNewPtr := reflect.New(srcValue.Type())
		dstNew := dstNewPtr.Elem()
		copyValue(dstNew, srcValue, all, parent)
		setValue(dst, dstNew, all)
		//dstP := *(*myValue)(unsafe.Pointer(&dst))
		//*(*interface{})(dstP.ptr) = dstNew.Interface()
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			copyValue(dst.Field(i), src.Field(i), all, parent)
		}
	case reflect.Chan:
		//todo 暂不拷贝
		//dstValuePtr := reflect.MakeChan(src.Type(), src.Cap())
		//setValue(dst, dstValuePtr)
	case reflect.Array:
		for i := 0; i < src.Len(); i++ {
			copyValue(dst.Index(i), src.Index(i), all, parent)
		}
	case reflect.Slice:
		if src.IsNil() {
			return
		}
		if !checkCanSet(dst, all) {
			return
		}
		dstNewSliceValue := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		for i := 0; i < src.Len(); i++ {
			copyValue(dstNewSliceValue.Index(i), src.Index(i), all, parent)
		}
		setValue(dst, dstNewSliceValue, all)

		//dstP := *(*myValue)(unsafe.Pointer(&dst))
		//srcP := *(*myValue)(unsafe.Pointer(&dstSliceValue))
		//*(*mySlice)(dstP.ptr) = *(*mySlice)(srcP.ptr)
	case reflect.String:
		if src.Len() == 0 {
			return
		}
		if !checkCanSet(dst, all) {
			return
		}
		setValue(dst, src, all)

		//dstP := *(*myValue)(unsafe.Pointer(&dst))
		//srcP := *(*myValue)(unsafe.Pointer(&src))
		//*(*myString)(dstP.ptr) = *(*myString)(srcP.ptr)
	case reflect.Map:
		if src.IsNil() {
			return
		}
		if !checkCanSet(dst, all) {
			return
		}
		dstNewMapValue := reflect.MakeMap(src.Type())
		for _, srcKey := range src.MapKeys() {
			srcValue := src.MapIndex(srcKey)
			dstNewValue := reflect.New(srcValue.Type()).Elem()
			copyValue(dstNewValue, srcValue, all, parent)
			dstNewKey := reflect.New(srcKey.Type()).Elem()
			copyValue(dstNewKey, srcKey, all, parent)
			dstNewMapValue.SetMapIndex(dstNewKey, dstNewValue)
		}
		setValue(dst, dstNewMapValue, all)

		//dstP := *(*myValue)(unsafe.Pointer(&dst))
		//srcP := *(*myValue)(unsafe.Pointer(&dstMapValue))
		//*(*unsafe.Pointer)(dstP.ptr) = srcP.ptr
	case reflect.Func:
		if src.IsNil() {
			return
		}
		if !checkCanSet(dst, all) {
			return
		}
		setValue(dst, src, all)
		//dstP := *(*myValue)(unsafe.Pointer(&dst))
		//srcP := *(*myValue)(unsafe.Pointer(&src))
		//*(*unsafe.Pointer)(dstP.ptr) = *(*unsafe.Pointer)(srcP.ptr)
	default:
		if !checkCanSet(dst, all) {
			return
		}
		log.ErrorF("==================%v", src)
		setValue(dst, src, all)
		//fmt.Printf("6666666666666: %v, %v, %v, %v\n", dst, src, dst.IsValid(), src.IsValid())
		//dstP := *(*myValue)(unsafe.Pointer(&dst))
		//srcP := *(*myValue)(unsafe.Pointer(&src))
		//*(*unsafe.Pointer)(dstP.ptr) = *(*unsafe.Pointer)(srcP.ptr)
	}
}
