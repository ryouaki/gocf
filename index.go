package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
#include "./invoke.h"
*/
import "C"
import "fmt"

var pluginMap = make(map[string]JSGoFuncHandler)

func InitGoCloudFunc(script string) {
	rt := NewRuntime()
	ctx := rt.NewContext()

	fb := NewJSGoFunc(ctx, func(args []*JSValue, this *JSValue) (*JSValue, *JSValue) {
		for _, v := range args {
			// fmt.Println(C.GoString(C.JS_ToCString(ctx.p, v.p)))
			val := v.GetPropertyByIndex(0)

			fmt.Println(C.GoString(C.JS_ToCString(val.ctx.p, val.p)))
		}
		fmt.Println("Invoke")
		return nil, nil
	})
	ctx.ExportFunc("test", fb)
	ret, _ := ctx.Eval(script, "main.js")
	fmt.Println(C.GoString(C.JS_ToCString(ctx.p, ret.p)))

	ret.Free()
	ctx.Free()
	rt.Free()

	// C.Invoke(nil, C.JSValue{}, C.int(0), nil)
}

func RegistPlugin(name string, fb JSGoFuncHandler) error {
	return nil
}
