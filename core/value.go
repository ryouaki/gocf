package core

/*
#include "./../quickjs-libc.h";
*/
import "C"

type JSValue struct {
	p   C.JSValue
	ctx *JSContext
}

func NewValue(ctx *JSContext, val C.JSValue) *JSValue {
	ret := new(JSValue)
	ret.ctx = ctx // 变量的上下文
	ret.p = val   // 变量在引擎中的值
	return ret
}
