package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
*/
import "C"

// 抽象的VM
type JSRuntime struct {
	p *C.JSRuntime // 指向JS引擎的指针
}

// 创建一个VM对象
func NewRuntime() *JSRuntime {
	rt := new(JSRuntime)
	rt.p = C.JS_NewRuntime()
	return rt
}

// 释放VM对象。
func (rt *JSRuntime) Free() {
	C.JS_FreeRuntime(rt.p)
}
