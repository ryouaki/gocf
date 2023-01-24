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
	"sync"
	"unsafe"
)

var vmLock sync.Mutex

type PluginCb struct {
	Name string
	Fb   JSGoFuncHandler
}

type JSVM struct {
	VM     *JSRuntime // 引擎
	Ctx    *JSContext // 执行上下文
	IsFree bool       // 是否被占用
}

var pluginMap = make(map[string][]*PluginCb)
var scripts = make([]string, 0, 8)
var vms = make([]JSVM, 0, 4)

func init() {
	nums := 1
	for idx, v := range os.Args {
		if v == "-n" && len(os.Args) > idx+1 {
			nums, _ = strconv.Atoi(os.Args[idx+1])
			break
		}
	}

	for i := 0; i < nums; i++ {
		rt := NewRuntime()
		ctx := rt.NewContext()
		vms = append(vms, JSVM{
			VM:     rt,
			Ctx:    ctx,
			IsFree: true,
		})
	}
}

func ReleaseVM(vm *JSVM) {
	vm.IsFree = true
}

func GetVM(timeout int, retry int) *JSVM {
	vmLock.Lock()
	defer vmLock.Unlock()
	for _, v := range vms {
		if v.IsFree {
			v.IsFree = false
			return &v
		}
	}
	return nil
}

func FreeVM() {
	for _, v := range vms {
		v.Ctx.Free()
		v.VM.Free()
	}
}

// 注册JS脚本，默认以文件路径为api请求地址
func RegistCloudFunc(path string) int {
	scripts = append(scripts, path)
	return len(scripts)
}

func InitGoCloudFunc() {
	for _, rt := range vms {
		for key, pls := range pluginMap {
			obj := NewObject(rt.Ctx)
			for _, fb := range pls {
				funcValue := &JSValue{
					Ctx: rt.Ctx,
					P:   NewJSGoFunc(rt.Ctx, fb.Fb).P,
				}
				obj.SetProperty(fb.Name, funcValue)
			}
			rt.Ctx.Global.SetProperty(key, obj)
		}
	}
}

func RunAPI(script string) {
	rt := GetVM(0, 0)

	_, err := rt.Ctx.Eval(script, "main", 1<<0|1<<5)
	if rt.Ctx.GetException() != nil {
		fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, err.P)))
		return
	}

	// C.js_module_set_import_meta(rt.Ctx.P, func_val.P, 0, 0)
	// module := C.GetModule(func_val.P)
	// if rt.Ctx.GetException() != nil {
	// 	r := rt.Ctx.GetException()
	// 	fmt.Println(r.ToString())
	// }

	cName := C.CString("main")
	defer C.free(unsafe.Pointer(cName))
	C.GetModule(rt.Ctx.P, C.JS_NewAtom(rt.Ctx.P, cName))
	// e := rt.Ctx.Global.GetProperty("exec")
	// fmt.Println(1, module)
	// C.FreeModule(rt.Ctx.P, module)
	if rt.Ctx.GetException() != nil {
		r := rt.Ctx.GetException()
		fmt.Println(r.ToString())
	}

	exec := `
	import exec from "main";

	console.log(333, exec)
	function exec1 (resolve, reject) {
		const _exec = exec
		_exec().then((res) => {
			console.log(res)
			resolve(res)
		}).catch(reject)
	};
	globalThis.exec1 = exec1
	`
	wfb, _ := rt.Ctx.Eval(exec, "", 1<<0)
	defer wfb.Free()
	if rt.Ctx.GetException() != nil {
		r := rt.Ctx.GetException()
		fmt.Println(r.ToString())
		module := C.GetModule(rt.Ctx.P, C.JS_NewAtom(rt.Ctx.P, cName))

		fmt.Println(module == nil)

	}
	// // C.ListModule(rt.Ctx.P)

	resolveCB := NewJSGoFunc(rt.Ctx, func(args []*JSValue, this *JSValue) *JSValue {
		fmt.Println(args[0].GetPropertyKeys().ToString())
		fmt.Println(args[0].GetProperty("data").GetProperty("a").ToString())
		return nil
	})

	rejectCB := NewJSGoFunc(rt.Ctx, func(args []*JSValue, this *JSValue) *JSValue {
		fmt.Println(args[0].ToString())
		return nil
	})

	callFb := rt.Ctx.Global.GetProperty("exec1")
	fmt.Println(callFb.IsFunction())

	args := []C.JSValue{
		resolveCB.P,
		rejectCB.P,
	}

	C.JS_Call(rt.Ctx.P, callFb.P, NewNull(rt.Ctx).P, 2, &args[0])
	// // fmt.Println(NewValue(rt.Ctx, result).IsException())
	if r := rt.Ctx.GetException(); r != nil {
		fmt.Println(r.ToString())
	}
	fmt.Println(C.JS_IsJobPending(rt.VM.P))
	if C.JS_IsJobPending(rt.VM.P) > 0 {
		C.JS_ExecutePendingJob(rt.VM.P, &rt.Ctx.P)
	}
	// if err != nil {
	// 	fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, err.P)))
	// } else {
	// 	// fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, ret.P)))
	// }
	ReleaseVM(rt)
}

// 注册JS可以调用的函数，挂载到global.gocf对象上
func RegistPlugin(name string, fbs []*PluginCb) error {
	_, had := pluginMap[name]
	if had {
		return fmt.Errorf("[GoCF]:Plugin \"%s\" has been registed.", name)
	}
	pluginMap[name] = fbs
	return nil
}
