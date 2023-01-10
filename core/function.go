package core

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "./../quickjs-libc.h";
*/
import "C"

type JSGoFuncHandler func(args []*JSValue, this *JSValue) (*JSValue, *JSValue)

type JSGoFunc struct {
	p   C.JSValue
	ctx *JSContext
	fb  JSGoFuncHandler
}

func NewJSGoFunc(ctx *JSContext, fb JSGoFuncHandler) *JSGoFunc {
	jsGoFunc := new(JSGoFunc)
	jsGoFunc.ctx = ctx
	jsGoFunc.fb = fb

	// 注入bridge
	ws := `(invoke, id) => function () {
		return invoke.call(this, id, arguments);
	}`

	// 这个执行后会返回一个函数的引用。
	wfb, _ := ctx.Eval(ws, "")
	defer wfb.Free()

	id := len(ctx.funcs)
	ctx.funcs = append(ctx.funcs, jsGoFunc) // 将自己加入到队列中

	cId := NewInt32(ctx, id)
	args := []C.JSValue{
		ctx.invokeFunc,
		cId.p,
	}

	jsGoFunc.p = C.JS_Call(ctx.p, wfb.p, NewNull(ctx).p, 2, &args[0])

	return jsGoFunc
}

func (fb *JSGoFunc) Export(name string) {
	fb.ctx.global.SetProperty(name, &JSValue{
		ctx: fb.ctx,
		p:   fb.p,
	})
}
