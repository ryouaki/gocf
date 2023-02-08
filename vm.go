package gocf

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L./ -lquickjs

#include "./quickjs-libc.h"
#include "./invoke.h"
*/
import "C"
import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

var vmLock sync.Mutex

// JS引擎结构体
type JSVM struct {
	VM     *JSRuntime // 引擎
	Ctx    *JSContext // 执行上下文
	IsFree bool       // 是否被占用
}

// JS引擎中可以使用的Go插件
type Plugin struct {
	Name string
	Cb   JSGoFuncHandler
}

// 引擎缓存池
var vms = make([]*JSVM, 0, 4)
var devVm *JSVM = nil

// Go插件缓存池
var pluginMap = make(map[string][]*Plugin)

func InitVM(nums int) {
	// 如果为0则直接退出
	if nums < 1 {
		return
	}

	// 根据启动参数-n初始化引擎数量
	for i := 0; i < nums; i++ {
		vms = append(vms, buildVM())
	}
	GoCFLog("initialized " + strconv.Itoa(nums) + " JSVM")
}

func InitDevVM() *JSVM {
	if devVm != nil {
		FreeDevVM()
	}
	devVm = buildVM()
	return devVm
}

func GetDevVM() *JSVM {
	return devVm
}

func FreeDevVM() {
	if devVm != nil {
		devVm.Ctx.Free()
		devVm.VM.Free()
	}
	devVm = nil
}

func buildVM() *JSVM {
	rt := NewRuntime()     // 初始化引擎
	ctx := rt.NewContext() // 初始化执行上下文

	// 初始化插件
	for key, pls := range pluginMap {
		obj := NewObject(ctx)
		for _, fb := range pls {
			funcValue := &JSValue{
				Ctx: ctx,
				P:   NewJSGoFunc(ctx, fb.Cb).P,
			}
			obj.SetProperty(fb.Name, funcValue)
		}
		ctx.Global.SetProperty(key, obj)
	}

	return &JSVM{
		VM:     rt,
		Ctx:    ctx,
		IsFree: true,
	}
}

// 释放虚拟机
func ReleaseVM(vm *JSVM) {
	vmLock.Lock()
	defer vmLock.Unlock()
	vm.IsFree = true
}

// 获取可用虚拟机
func GetVM(ot time.Duration) *JSVM {
	vmLock.Lock()         // 协程共享，需要加锁
	defer vmLock.Unlock() // 需要释放锁
	// 超时，不能一直查
	ctx, cancel := context.WithTimeout(context.Background(), ot*time.Millisecond)
	defer cancel()

	var vm *JSVM = nil
	rt := make(chan bool)

	go func() {
		for {
			for _, v := range vms {
				if v.IsFree {
					v.IsFree = false
					vm = v
					rt <- true
				}
			}
			// 不能一直占用cpu，需要休眠4ms
			time.Sleep(time.Duration(4) * time.Millisecond)
		}
	}()

	// 如果返回nil即为超时。没有获取到vm
	select {
	case <-rt:
		return vm
	case <-ctx.Done():
		return vm
	}
}

// 释放虚拟机
func FreeVM() {
	for _, v := range vms {
		v.Ctx.Free()
		v.VM.Free()
	}
}

// 注册JS可以调用的函数，挂载到global.gocf对象上
func RegistPlugin(name string, fbs []*Plugin) error {
	_, had := pluginMap[name]
	if had {
		return fmt.Errorf("Plugin \"%s\" has been registed.", name)
	}
	pluginMap[name] = fbs
	return nil
}

// 初始化api
func InitApi(isDev bool) error {
	var apis []ScriptApi
	if isDev {
		apis = ScriptDevApiMap
	} else {
		apis = ScriptApiMap
	}
	// 将Script脚本注入到各个VM的Ctx中。
	for _, v := range apis {
		// 获取完整脚本文件
		code, err := os.ReadFile(v.File)
		if err != nil {
			GoCFLog("Error", v.File+" Read failed", err.Error())
			return err
		}
		GoCFLog("Init API " + v.Method + " " + v.Module + " " + v.Path)
		if err = InjectModule(string(code), v, isDev); err != nil {
			GoCFLog("Error", v.File+" Eval failed", err.Error())
			return err
		}
	}
	return nil
}

func buildModule(m *JSVM, code string, name string) error {
	_, err := m.Ctx.Eval(code, name, 1<<0|1<<5)
	defer err.Free()
	if m.Ctx.GetException() != nil {
		return fmt.Errorf(err.ToString())
	}
	if err != nil {
		return fmt.Errorf(err.ToString())
	}
	return nil
}

func InjectModule(code string, api ScriptApi, isDev bool) error {
	if isDev {
		err := buildModule(devVm, code, api.Module)
		if err != nil {
			return err
		}
		GoCFLog("Inject module " + api.Module)
	} else {
		for _, v := range vms {
			err := buildModule(v, code, api.Module)
			if err != nil {
				return err
			}
			GoCFLog("Inject module " + api.Module)
		}
	}

	return nil
}

func FreeModule(name string) error {
	cStr := C.CString(name)
	defer C.free(unsafe.Pointer(cStr))
	for _, v := range vms {
		m := C.JS_FindLoadedModule(v.Ctx.P, C.JS_NewAtom(v.Ctx.P, cStr))
		if m != nil {
			C.JS_FreeModule(v.Ctx.P, m)
		}
	}
	return nil
}

func WaitForLoop(rt *JSVM) {
	for C.JS_IsJobPending(rt.VM.P) > 0 {
		C.JS_ExecutePendingJob(rt.VM.P, &rt.Ctx.P)
	}
}
