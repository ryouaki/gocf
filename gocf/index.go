package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
#include "./invoke.h"
*/
import "C"
import (
	"os"
	"strconv"
)

// 根据入参初始化参数
func init() {
	for idx, v := range os.Args {
		if v == "-n" && len(os.Args) > idx+1 {
			if nums, err := strconv.Atoi(os.Args[idx+1]); err != nil {
				Nums = nums
			}
		} else if v == "-p" && len(os.Args) > idx+1 {
			Root = os.Args[idx+1]
		}
	}
}

// 加载JS相关配置文件与引擎
func InitGoCloudFunc() {
	LoadApiScripts(Root)
	InitVM(Nums)
}

func RunAPI(script string) {
	// rt := GetVM(0, 0)

	// _, err := rt.Ctx.Eval(script, "main", 1<<0|1<<5)
	// if rt.Ctx.GetException() != nil {
	// 	fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, err.P)))
	// 	return
	// }

	// cStr := C.CString("tt/aa.js")
	// defer C.free(unsafe.Pointer(cStr))
	// m := C.JS_FindLoadedModule(rt.Ctx.P, C.JS_NewAtom(rt.Ctx.P, cStr))
	// fmt.Println(444, C.GoString(C.JS_ToCString(rt.Ctx.P, C.JS_AtomToString(rt.Ctx.P, C.JS_GetModuleName(rt.Ctx.P, m)))))

	// C.JS_FreeModule(rt.Ctx.P, m)

	// _, err = rt.Ctx.Eval(`
	// export default function () {
	// 	return "222"
	// }
	// `, "tt/aa.js", 1<<0|1<<5)
	// if rt.Ctx.GetException() != nil {
	// 	fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, err.P)))
	// 	return
	// }

	// resolveCB := NewJSGoFunc(rt.Ctx, func(args []*JSValue, this *JSValue) *JSValue {
	// 	fmt.Println(222, args[0].GetPropertyKeys().ToString())
	// 	fmt.Println(222, args[0].GetProperty("data").GetProperty("a").ToString())
	// 	return nil
	// })

	// rejectCB := NewJSGoFunc(rt.Ctx, func(args []*JSValue, this *JSValue) *JSValue {
	// 	fmt.Println(args[0].ToString())
	// 	return nil
	// })

	// rt.Ctx.Global.SetProperty("resolve", NewFunc(rt.Ctx, resolveCB))
	// rt.Ctx.Global.SetProperty("reject", NewFunc(rt.Ctx, rejectCB))

	// exec := `
	// import exec from "main";

	// exec().then((res) => {
	// 	console.log(res)
	// 	resolve(res)
	// }).catch(reject)
	// `
	// wfb, _ := rt.Ctx.Eval(exec, "<input>", 1<<0)
	// defer wfb.Free()
	// if rt.Ctx.GetException() != nil {
	// 	r := rt.Ctx.GetException()
	// 	fmt.Println(r.ToString())
	// }

	// // // // fmt.Println(NewValue(rt.Ctx, result).IsException())
	// if r := rt.Ctx.GetException(); r != nil {
	// 	fmt.Println(r.ToString())
	// }
	// fmt.Println(C.JS_IsJobPending(rt.VM.P))
	// if C.JS_IsJobPending(rt.VM.P) > 0 {
	// 	C.JS_ExecutePendingJob(rt.VM.P, &rt.Ctx.P)
	// }
	// if err != nil {
	// 	fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, err.P)))
	// } else {
	// 	// fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, ret.P)))
	// }
	// ReleaseVM(rt)
	// rt.Ctx.Free()
	// rt.VM.Free()
}
