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

	err, data := ctx.Funcs[id].Fb(args, &JSValue{
		Ctx: ctx,
		P:   cthis,
	})

	ret := NewArray(ctx)

	ret.SetPropertyByIndex(0, NewUndefined(ctx))
	ret.SetPropertyByIndex(1, NewUndefined(ctx))

	if err != nil {
		ret.SetPropertyByIndex(0, err)
	} else if data != nil {
		ret.SetPropertyByIndex(1, data)
	}

	return ret.P
}
