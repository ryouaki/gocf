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

type JSValue struct {
	P   C.JSValue
	Ctx *JSContext
}

func NewValue(ctx *JSContext, val C.JSValue) *JSValue {
	ret := new(JSValue)
	ret.Ctx = ctx // 变量的上下文
	ret.P = val   // 变量在引擎中的值
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

func (val *JSValue) GetPropertyKeys() *JSValue {
	keys := val.Ctx.Global.GetProperty("Object").GetProperty("keys")
	args := []C.JSValue{
		val.P,
	}
	ret := C.JS_Call(val.Ctx.P, keys.P, NewNull(val.Ctx).P, 1, &args[0])
	return NewValue(val.Ctx, ret)
}

func (val *JSValue) ToString() string {
	return C.GoString(C.JS_ToCString(val.Ctx.P, val.P))
}

func NewInt32(ctx *JSContext, d int) *JSValue {
	return &JSValue{
		Ctx: ctx,
		P:   C.JS_NewInt32(ctx.P, C.int32_t(int32(d))),
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

func NewArray(ctx *JSContext) *JSValue {
	return &JSValue{
		Ctx: ctx,
		P:   C.JS_NewArray(ctx.P),
	}
}

func (v *JSValue) IsNumber() bool {
	return C.JS_IsNumber(v.P) == 1
}
func (v *JSValue) IsBigInt() bool {
	return C.JS_IsBigInt(v.Ctx.P, v.P) == 1
}
func (v *JSValue) IsBigFloat() bool {
	return C.JS_IsBigFloat(v.P) == 1
}
func (v *JSValue) IsBigDecimal() bool {
	return C.JS_IsBigDecimal(v.P) == 1
}
func (v *JSValue) IsBool() bool {
	return C.JS_IsBool(v.P) == 1
}
func (v *JSValue) IsNull() bool {
	return C.JS_IsNull(v.P) == 1
}
func (v *JSValue) IsUndefined() bool {
	return C.JS_IsUndefined(v.P) == 1
}
func (v *JSValue) IsException() bool {
	return C.JS_IsException(v.P) == 1
}
func (v *JSValue) IsUninitialized() bool {
	return C.JS_IsUninitialized(v.P) == 1
}
func (v *JSValue) IsString() bool {
	return C.JS_IsString(v.P) == 1
}
func (v *JSValue) IsSymbol() bool {
	return C.JS_IsSymbol(v.P) == 1
}
func (v *JSValue) IsObject() bool {
	return C.JS_IsObject(v.P) == 1
}
func (v *JSValue) IsArray() bool {
	return C.JS_IsArray(v.Ctx.P, v.P) == 1
}
func (v *JSValue) IsError() bool {
	return C.JS_IsError(v.Ctx.P, v.P) == 1
}
func (v *JSValue) IsFunction() bool {
	return C.JS_IsFunction(v.Ctx.P, v.P) == 1
}
func (v *JSValue) IsConstructor() bool {
	return C.JS_IsConstructor(v.Ctx.P, v.P) == 1
}
