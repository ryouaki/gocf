package core

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lquickjs
#include "./invoke.h";
*/
import "C"
import "unsafe"

type JSContext struct {
	ctx        *C.JSContext
	funcs      []*JSGoFunc
	invokeFunc C.JSValue // 注入调用go函数的sdk api
	global     *JSValue
}

func (rt *JSRuntime) NewContext() *JSContext {
	ret := new(JSContext)
	ret.ctx = C.JS_NewContext(rt.p)
	ret.funcs = []*JSGoFunc{}                                                                              // 注册的Go API                                                                       // 引擎中的执行上下文句柄
	ret.invokeFunc = C.JS_NewCFunction(ret.ctx, (*C.JSCFunction)(unsafe.Pointer(C.Invoke)), nil, C.int(5)) // 调用go的api
	ret.global = NewValue(ret, C.JS_GetGlobalObject(ret.ctx))                                              // 全局对象句柄，用于挂载api
	return ret
}
