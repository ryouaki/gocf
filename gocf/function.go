package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "./quickjs-libc.h"
*/
import "C"

type JSGoFuncHandler func(args []*JSValue, this *JSValue) (*JSValue, *JSValue)

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
		var argvs = []
		for (var i = 0; i < arguments.length; i++) {
			var argv = arguments[i];
			if (typeof argv === "object") {
				argvs.push(JSON.stringify(argv))
			} else {
				argvs.push(argv)
			}
		}
		var [err, ret] = invoke.call(this, id, argvs);

		if (err !== undefined) {
			return {
				error: true,
				message: err
			}
		} else {
			return {
				error: false,
				message: ret
			}
		}
	}`

	// 这个执行后会返回一个函数的引用。
	wfb, _ := ctx.Eval(ws, "")
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
