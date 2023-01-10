package core

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "./../quickjs-libc.h";
*/
import "C"
import "unsafe"

const (
	IS_VALUE int = 0
	IS_ERROR int = -1
)

type JSValue struct {
	p   C.JSValue
	ctx *JSContext
	is  int
}

func NewValue(ctx *JSContext, val C.JSValue) *JSValue {
	ret := new(JSValue)
	ret.ctx = ctx // 变量的上下文
	ret.p = val   // 变量在引擎中的值
	ret.is = IS_VALUE
	return ret
}

func NewError(ctx *JSContext, val C.JSValue) *JSValue {
	ret := new(JSValue)
	ret.ctx = ctx // 变量的上下文
	ret.p = val   // 变量在引擎中的值
	ret.is = IS_ERROR
	return ret
}

func (val *JSValue) Free() {
	C.JS_FreeValue(val.ctx.p, val.p)
}

func (val *JSValue) SetProperty(key string, v *JSValue) {
	cStr := C.CString(key)
	defer C.free(unsafe.Pointer(cStr))
	C.JS_SetPropertyStr(val.ctx.p, val.p, cStr, v.p)
}

func (val *JSValue) GetProperty(key string) *JSValue {
	cStr := C.CString(key)
	defer C.free(unsafe.Pointer(cStr))
	return &JSValue{
		ctx: val.ctx,
		p:   C.JS_GetPropertyStr(val.ctx.p, val.p, cStr),
	}
}

func (val *JSValue) GetPropertyByIndex(idx int) *JSValue {
	return &JSValue{
		ctx: val.ctx,
		p:   C.JS_GetPropertyUint32(val.ctx.p, val.p, C.uint32_t(idx)),
	}
}

func (val *JSValue) IsError() bool {
	return val.is == IS_ERROR
}

func NewInt32(ctx *JSContext, d int) *JSValue {
	return &JSValue{
		ctx: ctx,
		p:   C.JS_NewInt32(ctx.p, C.int32_t(int32(d))),
		is:  IS_VALUE,
	}
}

func NewNull(ctx *JSContext) *JSValue {
	return &JSValue{
		ctx: ctx,
		p:   C.JS_NULL,
	}
}

func NewUndefined(ctx *JSContext) *JSValue {
	return &JSValue{
		ctx: ctx,
		p:   C.JS_UNDEFINED,
	}
}
