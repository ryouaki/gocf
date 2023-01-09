package core

/*
#include "./../quickjs-libc.h";
*/
import "C"

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

func (val *JSValue) IsError() bool {
	return val.is == IS_ERROR
}
