package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs
#include "./invoke.h"
*/
import "C"
import (
	"unsafe"
)

type JSContext struct {
	P          *C.JSContext
	Funcs      []*JSGoFunc
	InvokeFunc C.JSValue // 注入调用go函数的sdk api
	Global     *JSValue
}

type JSGoFuncHandler func(args []*JSValue, this *JSValue) *JSValue

type JSGoFunc struct {
	P   C.JSValue
	Ctx *JSContext
	Fb  JSGoFuncHandler
}

var ctxCache = make(map[*C.JSContext]*JSContext)

func (rt *JSRuntime) NewContext() *JSContext {
	ret := new(JSContext)
	ret.P = C.NewJsContext(rt.P)
	// 引擎中的执行上下文句柄
	ret.Funcs = []*JSGoFunc{}
	// 调用go的api
	ret.InvokeFunc = C.JS_NewCFunction(ret.P, (*C.JSCFunction)(unsafe.Pointer(C.Invoke)), nil, C.int(5))
	ret.Global = NewValue(ret, C.JS_GetGlobalObject(ret.P)) // 全局对象句柄，用于挂载api

	ctxCache[ret.P] = ret
	return ret
}

func (ctx *JSContext) Eval(script string, filename string, flag int) (*JSValue, *JSValue) {
	jsStr := C.CString(script)          // 将JS文本代码转换为quickjs引擎代码格式
	defer C.free(unsafe.Pointer(jsStr)) // 执行结束后需要释放空间

	jsStrLen := C.ulong(len(script))         // js代码长度
	jsFileName := C.CString(filename)        // js文件名。quickjs每次eval都要文件名
	defer C.free(unsafe.Pointer(jsFileName)) // 执行结束后释放空间

	ret := &JSValue{
		P:   C.JS_Eval(ctx.P, jsStr, jsStrLen, jsFileName, C.int(flag)),
		Ctx: ctx,
	}

	err := ctx.GetException()
	defer err.Free()

	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (ctx *JSContext) GetException() *JSValue {
	err := C.JS_GetException(ctx.P)
	if C.JS_IsNull(err) == 1 {
		return nil
	}
	return &JSValue{
		Ctx: ctx,
		P:   err,
	}
}

func (ctx *JSContext) ExportValue(name string, val *JSValue) {
	ctx.Global.SetProperty(name, val)
}

func (ctx *JSContext) ExportFunc(name string, fb *JSGoFunc) {
	ctx.Global.SetProperty(name, &JSValue{
		Ctx: fb.Ctx,
		P:   fb.P,
	})
}

// 释放JS引擎内变量空间
func (ctx *JSContext) FreeJSValue(val C.JSValue) {
	C.JS_FreeValue(ctx.P, val)
}

// 释放Go层映射JS变量
func (ctx *JSContext) FreeValue(val *JSValue) {
	C.JS_FreeValue(ctx.P, val.P)
}

// 释放Ctx
func (ctx *JSContext) Free() {
	// clean plugins
	ctx.Funcs = ctx.Funcs[0:0]
	for key, pls := range pluginMap {
		root := ctx.Global.GetProperty(key)
		for _, fb := range pls {
			fb.p.Free()
			// root.DeleteProperty(fb.Name)
		}
		root.Free()
		// ctx.Global.DeleteProperty(key)
	}

	ctx.FreeJSValue(ctx.InvokeFunc)
	// fmt.Println(ctx.Global.GetPropertyKeys().ToString())
	ctx.Global.Free()
	_, key := ctxCache[ctx.P]
	if key {
		delete(ctxCache, ctx.P)
	}
	C.JS_FreeContext(ctx.P)
}

// 查询模块
func (ctx *JSContext) FindModule(name string) bool {
	cStr := C.CString(name)
	defer C.free(unsafe.Pointer(cStr))
	m := C.JS_FindLoadedModule(ctx.P, C.JS_NewAtom(ctx.P, cStr))
	if m != nil {
		return true
	} else {
		return false
	}
}

func NewJSGoFunc(ctx *JSContext, fb JSGoFuncHandler) *JSGoFunc {
	jsGoFunc := new(JSGoFunc)
	jsGoFunc.Ctx = ctx
	jsGoFunc.Fb = fb

	invokeFunc := ctx.Global.GetProperty("$$invoke")
	if invokeFunc == nil || !invokeFunc.IsFunction() {
		// 注入bridge
		ws := `globalThis.$$invoke = function (invoke, id) {
			return function () {
				var argvs = [id]
				for (var i = 0; i < arguments.length; i++) {
					var argv = arguments[i];
					argvs.push(argv)
				}
				var ret = invoke.apply(this, argvs);
				try {
					objData = JSON.parse(ret.data)
					ret.data = objData
				} catch(e) {}
				return ret
			}
		};`

		// 这个执行后会返回一个函数的引用。
		wfb, e := ctx.Eval(ws, "<code>", 1<<0)
		defer wfb.Free()
		defer e.Free()

		if e != nil {
			GoCFLog(e.ToString())
		}

		r := ctx.GetException()
		if r != nil {
			GoCFLog(r.ToString())
			r.Free()
		}

		invokeFunc = ctx.Global.GetProperty("$$invoke")
	}
	// defer invokeFunc.Free()

	id := len(ctx.Funcs)
	ctx.Funcs = append(ctx.Funcs, jsGoFunc) // 将自己加入到队列中

	cId := NewInt32(ctx, id)
	defer cId.Free()
	args := []C.JSValue{
		ctx.InvokeFunc,
		cId.P,
	}

	jsGoFunc.P = C.JS_Call(ctx.P, invokeFunc.P, C.JS_NULL, 2, &args[0])

	return jsGoFunc
}

func MakeInvokeResult(ctx *JSContext, status int, val interface{}) *JSValue {
	data := InterfaceToString(val)
	cStr := C.CString(data)
	defer C.free(unsafe.Pointer(cStr))
	cVal := C.JS_NewString(ctx.P, cStr)
	ret := NewObject(ctx)
	ret.SetProperty("error", NewValue(ctx, C.JS_NewBool(ctx.P, C.int(status))))
	ret.SetProperty("data", NewValue(ctx, cVal))

	return ret
}
