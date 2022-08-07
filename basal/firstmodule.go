package basal

import (
	"github.com/jingyanbin/core/internal"
	"reflect"
	"runtime"
	"unsafe"
)

type functab struct {
	entry   uintptr
	funcoff uintptr
}

// Mapping information for secondary text sections

type textsect struct {
	vaddr    uintptr // prelinked section vaddr
	length   uintptr // section length
	baseaddr uintptr // relocated section address
}

type itab struct {
	inter interface{}
	_type interface{}
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type pcHeader struct {
	magic          uint32  // 0xFFFFFFFA
	pad1, pad2     uint8   // 0,0
	minLC          uint8   // min instruction size
	ptrSize        uint8   // size of a ptr in bytes
	nfunc          int     // number of functions in the module
	nfiles         uint    // number of entries in the file tab.
	funcnameOffset uintptr // offset to the funcnametab variable from pcHeader
	cuOffset       uintptr // offset to the cutab variable from pcHeader
	filetabOffset  uintptr // offset to the filetab variable from pcHeader
	pctabOffset    uintptr // offset to the pctab varible from pcHeader
	pclnOffset     uintptr // offset to the pclntab variable from pcHeader
}

type nameOff int32
type typeOff int32
type textOff int32

type ptabEntry struct {
	name nameOff
	typ  typeOff
}

type moduledata struct {
	pcHeader     *pcHeader
	funcnametab  []byte
	cutab        []uint32
	filetab      []byte
	pctab        []byte
	pclntable    []byte
	ftab         []functab
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr

	textsectmap []textsect
	typelinks   []int32 // offsets from types
	itablinks   []*itab

	ptab []ptabEntry

	next *moduledata
}

// //go:linkname firstmoduledata runtime.firstmoduledata
var firstmoduledata moduledata

// //go:linkname lastmoduledatap runtime.lastmoduledatap
var lastmoduledatap *moduledata // linker symbol

// //go:linkname modulesSlice runtime.modulesSlice
var modulesSlice *[]*moduledata // see activeModules

func FindFuncWithName(name string) (uintptr, error) {
	for moduleData := &firstmoduledata; moduleData != nil; moduleData = moduleData.next {
		for _, ftab := range moduleData.ftab {
			myFunc := runtime.FuncForPC(ftab.entry)
			if myFunc == nil {
				continue
			}
			//fmt.Println(myFunc.Name())
			if myFunc.Name() == name {
				return myFunc.Entry(), nil
			}
		}
	}
	//
	//for moduleData := lastmoduledatap; moduleData != nil; moduleData = moduleData.next {
	//	for _, ftab := range moduleData.ftab {
	//		myFunc := runtime.FuncForPC(ftab.entry)
	//		if myFunc == nil {
	//			continue
	//		}
	//		fmt.Println(myFunc.Name())
	//		if myFunc.Name() == name {
	//			return myFunc.Entry(), nil
	//		}
	//	}
	//}
	//
	//for _, moduleData := range *modulesSlice {
	//	for _, ftab := range moduleData.ftab {
	//		myFunc := runtime.FuncForPC(ftab.entry)
	//		if myFunc == nil {
	//			continue
	//		}
	//		fmt.Println(myFunc.Name())
	//		if myFunc.Name() == name {
	//			return myFunc.Entry(), nil
	//		}
	//	}
	//}
	return 0, NewError("invalid function " + name)
}

func GetFunc(outFuncPtr interface{}, name string) (err error) {
	defer internal.CatchError(func(e error) {
		err = NewError("exception: %v", e)
	})
	if IsPointer(outFuncPtr) == false {
		return NewError("not is func ptr")
	}

	var codePtr uintptr
	codePtr, err = FindFuncWithName(name)
	if err == nil {
		err = CreateFuncForCodePtr(outFuncPtr, codePtr)
	}
	return
}

type function struct {
	codePtr uintptr
}

func CreateFuncForCodePtr(outFuncPtr interface{}, codePtr uintptr) (err error) {
	defer internal.CatchError(func(e error) {
		err = e
	})
	outFuncVal := reflect.ValueOf(outFuncPtr).Elem()
	newFuncVal := reflect.MakeFunc(outFuncVal.Type(), nil)
	funcValuePtr := reflect.ValueOf(newFuncVal).FieldByName("ptr").Pointer()
	funcPtr := (*function)(unsafe.Pointer(funcValuePtr))
	funcPtr.codePtr = codePtr
	outFuncVal.Set(newFuncVal)
	return
}

type _func struct {
	entry   uintptr // start pc
	nameoff int32   // function name

	args   int32  // in/out args size
	funcID uint32 // set for certain special runtime functions

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	nfuncdata int32
}
