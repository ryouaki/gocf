package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "./quickjs-libc.h"
*/
import "C"
import (
	"unsafe"
)

const (
	CB_SUCCESS = 1
	CB_FAILED  = 0
)

type JSGoFuncHandler func(args []*JSValue, this *JSValue) *JSValue

type JSGoFunc struct {
	P   C.JSValue
	Ctx *JSContext
	Fb  JSGoFuncHandler
}

func NewJSGoFunc(ctx *JSContext, fb JSGoFuncHandler) *JSGoFunc {
	jsGoFunc := new(JSGoFunc)
	jsGoFunc.Ctx = ctx
	jsGoFunc.Fb = fb

	// 注入bridge
	ws := `(invoke, id) => function () {
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
	}`

	// 这个执行后会返回一个函数的引用。
	wfb, _ := ctx.Eval(ws, "", 0)
	defer wfb.Free()

	id := len(ctx.Funcs)
	ctx.Funcs = append(ctx.Funcs, jsGoFunc) // 将自己加入到队列中

	cId := NewInt32(ctx, id)
	args := []C.JSValue{
		ctx.InvokeFunc,
		cId.P,
	}

	jsGoFunc.P = C.JS_Call(ctx.P, wfb.P, NewNull(ctx).P, 2, &args[0])

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
