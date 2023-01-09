package core

/*
#include "./../quickjs-libc.h";
#include "./invoke.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	IS_VALUE int = 0
	IS_ERROR int = -1
)

type JSRuntime struct {
	cRuntime *C.JSRuntime
}

type JSContext struct {
	cContext *C.JSContext
	global   *JSValue
}

type JSValue struct {
	cValue C.JSValue
	ctx    *JSContext
	is     int
}

type JSError struct {
	cValue C.JSValue
	ctx    *JSContext
	val    *JSValue
}

func newJSRuntime() *JSRuntime {
	ret := &JSRuntime{
		cRuntime: C.JS_NewRuntime(),
	}

	return ret
}

func (runtime *JSRuntime) newContext() *JSContext {
	ctx := &JSContext{
		cContext: C.JS_NewContext(runtime.cRuntime),
	}
	ctx.global = &JSValue{
		ctx:    ctx,
		cValue: C.JS_GetGlobalObject(ctx.cContext),
	}

	return ctx
}

func (context *JSContext) Eval(script, filename string) (*JSValue, *JSError) {
	cScript := C.CString(script) // 转化代码到Quickjs的字符串类型
	defer C.free(unsafe.Pointer(cScript))
	length := C.ulong(len(script)) // 脚本原长度。JS_Eval入参需要

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	ret := &JSValue{
		ctx:    context,
		cValue: C.JS_Eval(context.cContext, cScript, length, cFilename, C.int(0)),
	}
	cstr := C.JS_ToPropertyKey(ret.ctx.cContext, ret.cValue)
	fmt.Println(cstr)
	// err := context.Exception()

	return nil, nil
}

// func (context *JSContext) Exception() *JSError {
// 	val := &JSValue{
// 		ctx:    context,
// 		cValue: C.JS_GetException(context.cContext),
// 	}

// 	return val
// }

func (runtime *JSRuntime) free() {
	C.JS_FreeRuntime(runtime.cRuntime)
}

func (context *JSContext) free() {
	C.JS_FreeContext(context.cContext)
}

func InitGoCloudFunc(script string) {
	rt := newJSRuntime()
	ctx := rt.newContext()

	ctx.Eval(script, "main.js")

	ctx.free()
	rt.free()
	// fmt.Println(ret, err)
	// C.Invoke(nil, C.JSValue{}, C.int(0), nil)
}
