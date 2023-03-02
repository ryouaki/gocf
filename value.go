package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type JSValue struct {
	P   C.JSValue
	Ctx *JSContext
}

type JSAtom struct {
	P   C.JSAtom
	Ctx *JSContext
}

func NewValue(ctx *JSContext, val C.JSValue) *JSValue {
	ret := new(JSValue)
	ret.Ctx = ctx // 变量的上下文
	ret.P = val   // 变量在引擎中的值
	return ret
}

func (val *JSValue) Free() {
	if val != nil {
		C.JS_FreeValue(val.Ctx.P, val.P)
	}
}

func (atom *JSAtom) Free() {
	C.JS_FreeAtom(atom.Ctx.P, atom.P)
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

func (val *JSValue) DeleteProperty(key string) error {
	v := NewString(val.Ctx, key)
	defer v.Free()
	a := C.JS_ValueToAtom(val.Ctx.P, v.P)
	defer C.JS_FreeAtom(val.Ctx.P, a)
	e := val.Ctx.GetException()
	defer e.Free()
	res := C.JS_DeleteProperty(val.Ctx.P, val.P, a, 0)

	if res < 0 {
		return fmt.Errorf("Delete " + key + " failed")
	}

	return nil
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

func (val *JSValue) GetPropertyKeys() []*JSValue {
	var (
		ptr  *C.JSPropertyEnum
		size C.uint32_t
	)

	result := int(C.JS_GetOwnPropertyNames(val.Ctx.P, &ptr, &size, val.P, C.int(1<<0|1<<1|1<<2)))
	if result < 0 {
		return []*JSValue{}
	}
	defer C.js_free(val.Ctx.P, unsafe.Pointer(ptr))

	entries := (*[1 << unsafe.Sizeof(0)]C.JSPropertyEnum)(unsafe.Pointer(ptr))

	names := make([]*JSValue, size)

	for i := 0; C.uint32_t(i) < size; i++ {
		v := NewValue(val.Ctx, C.JS_AtomToValue(val.Ctx.P, entries[i].atom))
		names[i] = v
		v.Free()
		C.JS_FreeAtom(val.Ctx.P, entries[i].atom)
	}

	return names
}

func (val *JSValue) ToString() string {
	if val == nil {
		return ""
	}
	return C.GoString(C.JS_ToCString(val.Ctx.P, val.P))
}

func NewString(ctx *JSContext, key string) *JSValue {
	cStr := C.CString(key)
	defer C.free(unsafe.Pointer(cStr))

	return &JSValue{
		Ctx: ctx,
		P:   C.JS_NewString(ctx.P, cStr),
	}
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

func NewFunc(ctx *JSContext, f *JSGoFunc) *JSValue {
	return &JSValue{
		Ctx: ctx,
		P:   f.P,
	}
}

func NewAtom(ctx *JSContext, v *JSValue) *JSAtom {
	return &JSAtom{
		Ctx: ctx,
		P:   C.JS_ValueToAtom(ctx.P, v.P),
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
