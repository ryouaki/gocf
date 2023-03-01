package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
*/
import "C"
import "fmt"

// 抽象的VM
type JSRuntime struct {
	P *C.JSRuntime // 指向JS引擎的指针
}

// 创建一个VM对象
func NewRuntime() *JSRuntime {
	rt := new(JSRuntime)
	rt.P = C.JS_NewRuntime()
	return rt
}

// 释放VM对象。
func (rt *JSRuntime) Free() {
	defer func() {
		if err := recover(); err != nil { //注意必须要判断
			fmt.Println(err)
		}
	}() //用来调用此匿名函数
	C.JS_FreeRuntime(rt.P)
}
