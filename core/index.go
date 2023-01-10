package core

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "./../quickjs-libc.h";
#include "./invoke.h"
*/
import "C"
import "fmt"

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
	fb.Export("test")
	ret, _ := ctx.Eval(script, "main.js")
	// fmt.Println(C.GoString(C.JS_ToCString(ctx.p, err.p)), err)

	ret.Free()
	ctx.Free()
	rt.Free()

	// C.Invoke(nil, C.JSValue{}, C.int(0), nil)
}
