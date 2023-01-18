package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
*/
import "C"
import "unsafe"

const (
	IS_VALUE int = 0
	IS_ERROR int = -1
)

type JSValue struct {
	P   C.JSValue
	Ctx *JSContext
	Is  int
}

func NewValue(ctx *JSContext, val C.JSValue) *JSValue {
	ret := new(JSValue)
	ret.Ctx = ctx // 变量的上下文
	ret.P = val   // 变量在引擎中的值
	ret.Is = IS_VALUE
	return ret
}

func NewError(ctx *JSContext, val C.JSValue) *JSValue {
	ret := new(JSValue)
	ret.Ctx = ctx // 变量的上下文
	ret.P = val   // 变量在引擎中的值
	ret.Is = IS_ERROR
	return ret
}

func (val *JSValue) Free() {
	C.JS_FreeValue(val.Ctx.P, val.P)
}

func (val *JSValue) SetProperty(key string, v *JSValue) {
	cStr := C.CString(key)
	defer C.free(unsafe.Pointer(cStr))
	C.JS_SetPropertyStr(val.Ctx.P, val.P, cStr, v.P)
}

func (val *JSValue) GetProperty(key string) *JSValue {
	cStr := C.CString(key)
	defer C.free(unsafe.Pointer(cStr))
	return &JSValue{
		Ctx: val.Ctx,
		P:   C.JS_GetPropertyStr(val.Ctx.P, val.P, cStr),
	}
}

func (val *JSValue) SetPropertyByIndex(idx int, v *JSValue) {
	C.JS_SetPropertyInt64(val.Ctx.P, val.P, C.int64_t(idx), v.P)
}

func (val *JSValue) GetPropertyByIndex(idx int) *JSValue {
	return &JSValue{
		Ctx: val.Ctx,
		P:   C.JS_GetPropertyUint32(val.Ctx.P, val.P, C.uint32_t(idx)),
	}
}

func (val *JSValue) IsError() bool {
	return val.Is == IS_ERROR
}

func (val *JSValue) ToString() string {
	return C.GoString(C.JS_ToCString(val.Ctx.P, val.P))
}

func NewInt32(ctx *JSContext, d int) *JSValue {
	return &JSValue{
		Ctx: ctx,
		P:   C.JS_NewInt32(ctx.P, C.int32_t(int32(d))),
		Is:  IS_VALUE,
	}
}

func NewNull(ctx *JSContext) *JSValue {
	return &JSValue{
		Ctx: ctx,
		P:   C.JS_NULL,
	}
}

func NewUndefined(ctx *JSContext) *JSValue {
	return &JSValue{
		Ctx: ctx,
		P:   C.JS_UNDEFINED,
	}
}

func NewObject(ctx *JSContext) *JSValue {
	return &JSValue{
		Ctx: ctx,
		P:   C.JS_NewObject(ctx.P),
	}
}
