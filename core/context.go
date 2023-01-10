package core

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs
#include "./invoke.h";
*/
import "C"
import "unsafe"

var ctxCache = make(map[*C.JSContext]*JSContext)

type JSContext struct {
	p          *C.JSContext
	funcs      []*JSGoFunc
	invokeFunc C.JSValue // 注入调用go函数的sdk api
	global     *JSValue
}

func (rt *JSRuntime) NewContext() *JSContext {
	ret := new(JSContext)
	ret.p = C.JS_NewContext(rt.p)
	// 引擎中的执行上下文句柄
	ret.funcs = []*JSGoFunc{}
	// 调用go的api
	ret.invokeFunc = C.JS_NewCFunction(ret.p, (*C.JSCFunction)(unsafe.Pointer(C.Invoke)), nil, C.int(5))
	ret.global = NewValue(ret, C.JS_GetGlobalObject(ret.p)) // 全局对象句柄，用于挂载api

	ctxCache[ret.p] = ret
	return ret
}

func (ctx *JSContext) Eval(script string, filename string) (*JSValue, *JSValue) {
	jsStr := C.CString(script)          // 将JS文本代码转换为quickjs引擎代码格式
	defer C.free(unsafe.Pointer(jsStr)) // 执行结束后需要释放空间

	jsStrLen := C.ulong(len(script))         // js代码长度
	jsFileName := C.CString(filename)        // js文件名。quickjs每次eval都要文件名
	defer C.free(unsafe.Pointer(jsFileName)) // 执行结束后释放空间

	ret := &JSValue{
		p:   C.JS_Eval(ctx.p, jsStr, jsStrLen, jsFileName, C.int(0)),
		ctx: ctx,
	}

	err := ctx.GetException()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (ctx *JSContext) GetException() *JSValue {
	err := C.JS_GetException(ctx.p)
	if C.JS_IsNull(err) == 1 {
		return nil
	}
	return &JSValue{
		ctx: ctx,
		p:   err,
		is:  IS_ERROR,
	}
}

// 释放JS引擎内变量空间
func (ctx *JSContext) FreeJSValue(val C.JSValue) {
	C.JS_FreeValue(ctx.p, val)
}

// 释放Go层映射JS变量
func (ctx *JSContext) FreeValue(val *JSValue) {
	C.JS_FreeValue(ctx.p, val.p)
}

// 释放Ctx
func (ctx *JSContext) Free() {
	ctx.FreeJSValue(ctx.invokeFunc)
	ctx.global.Free()
	_, key := ctxCache[ctx.p]
	if key {
		delete(ctxCache, ctx.p)
	}
	C.JS_FreeContext(ctx.p)
}
