package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
*/
import "C"
import (
	"unsafe"
)

//export GoInvoke
func GoInvoke(cctx *C.JSContext, cthis C.JSValueConst, cargc C.int, cargv *C.JSValueConst) C.JSValue {
	ctx := ctxCache[cctx]
	cArgs := (*[1 << 28]C.JSValueConst)(unsafe.Pointer(cargv))[:cargc:cargc]
	i := C.int64_t(0)
	C.JS_ToInt64(cctx, &i, cArgs[0])
	id := int(i)
	var args []*JSValue
	for _, cArg := range cArgs[1:] {
		args = append(args, &JSValue{
			Ctx: ctx,
			P:   cArg,
		})
	}

	ret, err := ctx.Funcs[id].Fb(args, &JSValue{
		Ctx: ctx,
		P:   cthis,
	})
	if err != nil {
		return err.P
	} else if ret == nil {
		return C.JS_UNDEFINED
	}

	return ret.P
}
