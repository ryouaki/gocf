package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
#include "./invoke.h"
*/
import "C"
import (
	"fmt"
	"os"
	"strconv"
)

// 启动引擎数量，默认1
var Nums = 1

// 脚本加载目录，默认./
var Root = "./"

// Master的地址+端口
var MasterHost = "127.0.0.1:8001"

// 根据入参初始化参数
func init() {
	for idx, v := range os.Args {
		if v == "-n" && len(os.Args) > idx+1 {
			if nums, err := strconv.Atoi(os.Args[idx+1]); err != nil {
				Nums = nums
			}
		} else if v == "-p" && len(os.Args) > idx+1 {
			Root = os.Args[idx+1]
		} else if v == "-m" && len(os.Args) > idx+1 {
			MasterHost = os.Args[idx+1]
		}
	}
}

// 加载JS相关配置文件与引擎
func InitGoCloudFunc() {
	LoadApiScripts(Root+"/api", false, "/api")
	InitVM(Nums)
	InitApi(false)
}

func RunAPI() {
	rt := GetVM(1)

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

	exec := `
	console.log(111)
	`
	wfb, err := rt.Ctx.Eval(exec, "<input>", 1<<0)
	defer wfb.Free()
	if rt.Ctx.GetException() != nil {
		r := rt.Ctx.GetException()
		fmt.Println(r.ToString())
	}

	if err != nil {
		fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, err.P)))
	}
	// ReleaseVM(rt)
	// rt.Ctx.Free()
	// rt.VM.Free()
}
