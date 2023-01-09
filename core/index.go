package core

/*
#include "./../quickjs-libc.h";
#include "./invoke.h"
*/
import "C"
import "fmt"

func InitGoCloudFunc(script string) {
	rt := NewRuntime()
	ctx := rt.NewContext()

	ret, err := ctx.Eval(script, "main.js")
	fmt.Println(C.GoString(C.JS_ToCString(ctx.p, ret.p)), err)
	ret.Free()
	ctx.Free()
	rt.Free()

	// C.Invoke(nil, C.JSValue{}, C.int(0), nil)
}
