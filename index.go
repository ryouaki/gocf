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
var scripts = make([]string, 4, 8)
var vms = make([]JSVM, 2, 4)

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

func InitGoCloudFunc(script string) {
	// rt := NewRuntime()
	// ctx := rt.NewContext()

	// fb := NewJSGoFunc(ctx, func(args []*JSValue, this *JSValue) (*JSValue, *JSValue) {
	// 	for _, v := range args {
	// 		// fmt.Println(C.GoString(C.JS_ToCString(ctx.p, v.p)))
	// 		val := v.GetPropertyByIndex(0)

	// 		fmt.Println(C.GoString(C.JS_ToCString(val.ctx.p, val.p)))
	// 	}
	// 	fmt.Println("Invoke")
	// 	return nil, nil
	// })
	// ctx.ExportFunc("test", fb)
	// ret, _ := ctx.Eval(script, "main.js")
	// fmt.Println(C.GoString(C.JS_ToCString(ctx.p, ret.p)))

	// ret.Free()
	// ctx.Free()
	// rt.Free()

	// C.Invoke(nil, C.JSValue{}, C.int(0), nil)
}

// 注册JS可以调用的函数，挂载到global.gocf对象上
func RegistPlugin(name string, fbs []*PluginCb) error {
	_, had := pluginMap[name]
	if had {
		return fmt.Errorf("[GoCF] Plugin \"%s\" has been registed.", name)
	}
	pluginMap[name] = fbs
	return nil
}
