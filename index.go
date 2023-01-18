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
	ret, _ := rt.Ctx.Eval(script, "main.js")
	fmt.Println(C.GoString(C.JS_ToCString(rt.Ctx.P, ret.P)))
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
