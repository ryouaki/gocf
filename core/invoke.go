package core

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "../quickjs-libc.h"
*/
import "C"
import "fmt"

//export GoInvoke
func GoInvoke(cctx *C.JSContext, cthis C.JSValueConst, cargc C.int, cargv *C.JSValueConst) C.JSValue {
	val := C.JSValue{}
	fmt.Println("Success")
	return val
}
