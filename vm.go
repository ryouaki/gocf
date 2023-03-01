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
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

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

type ScriptApi struct {
	Path   string // api地址
	Module string // 模块名
	Method string // api方法
	File   string // 脚本文件地址
}

const (
	METHOD_GET     = "get"
	METHOD_POST    = "post"
	METHOD_OPTIONS = "options"
	METHOD_PUT     = "put"
	METHOD_DELETE  = "delete"
	METHOD_PATCH   = "patch"
	METHOD_HEAD    = "head"
)

var methods = []string{
	METHOD_GET,
	METHOD_POST,
	METHOD_OPTIONS,
	METHOD_PUT,
	METHOD_DELETE,
	METHOD_PATCH,
	METHOD_HEAD,
}

// 引擎缓存池的互斥锁
var vmLock sync.Mutex
var resetLock sync.Mutex

// 启动引擎数量，默认1
var Nums = 1

// 脚本加载目录，默认./
var Root = "./"

// 引擎缓存池
var vms = make([]*JSVM, 0, 2)

// Go插件缓存池
var pluginMap = make(map[string][]*Plugin)

// api mapping
var ScriptApiMap = make([]ScriptApi, 0, 4)

var MasterHost = "http://localhost:8000"

var isReseting = false

// 根据入参初始化参数
func RunGoCF() {
	for idx, v := range os.Args {
		if v == "-n" && len(os.Args) > idx+1 {
			if nums, err := strconv.Atoi(os.Args[idx+1]); err != nil {
				Nums = nums
			}
		} else if v == "-p" && len(os.Args) > idx+1 {
			Root = os.Args[idx+1]
		} else if v == "-h" && len(os.Args) > idx+1 {
			MasterHost = os.Args[idx+1]
		}
	}

	InitVM(Nums)
	LoadApiScripts(Root+"/api", "/api")
	InitApi()
}

// ******************** VM 初始化 Start *************************
func InitVM(nums int) {
	// 如果为0则直接退出
	if nums < 1 {
		return
	}

	// 根据启动参数-n初始化引擎数量
	for i := 0; i < nums; i++ {
		vms = append(vms, newVM())
	}
	GoCFLog("initialized " + strconv.Itoa(nums) + " JSVM")
}

func newVM() *JSVM {
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
	if isReseting {
		return nil
	}
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
					return
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

func WaitForLoop(rt *JSVM) {
	for C.JS_IsJobPending(rt.VM.P) > 0 {
		C.JS_ExecutePendingJob(rt.VM.P, &rt.Ctx.P)
	}
}

func updateResetFlag(reset bool) {
	resetLock.Lock()
	defer resetLock.Unlock()
	isReseting = reset
}

// ******************** VM 初始化 End *************************

// ******************** Go插件 初始化 Start *************************
// 注册JS可以调用的函数，挂载到global.gocf对象上
func RegistPlugin(name string, fbs []*Plugin) error {
	_, had := pluginMap[name]
	if had {
		GoCFLog("Plugin \"%s\" has been registed.", name)
	}
	pluginMap[name] = fbs
	return nil
}

// ******************** Go插件 初始化 End *************************

// ******************** Api 初始化 Start *************************
func InitApi() error {
	apis := ScriptApiMap

	// 将Script脚本注入到各个VM的Ctx中。
	for _, v := range apis {
		// 获取完整脚本文件
		code, err := os.ReadFile(v.File)
		if err != nil {
			GoCFLog("Error", v.File+" Read failed", err.Error())
			return err
		}
		GoCFLog("Init API " + v.Method + " " + v.Module + " " + v.Path)
		if err = InjectModule(string(code), v); err != nil {
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

func InjectModule(code string, api ScriptApi) error {
	for _, v := range vms {
		err := buildModule(v, code, api.Module)
		if err != nil {
			return err
		}
		GoCFLog("Inject module " + api.Module)
	}

	return nil
}

// 清空 api
func ClearApiMap() {
	ScriptApiMap = ScriptApiMap[0:0]
}

// 根据参数指定目录加载脚本
func LoadApiScripts(root string, parent string) error {
	GoCFLog(root)

	// 遍历脚本文件目录
	dir, err := ioutil.ReadDir(root)
	if err != nil {
		GoCFLog("Error", "Load Script file failed!")
		return err
	}

	for _, f := range dir {
		name := f.Name()
		if f.IsDir() {
			LoadApiScripts(root+"/"+name, parent+"/"+name)
		} else if !strings.HasSuffix(name, "js") {
			continue
		} else {
			apiInfo := strings.Split(name, ".")
			if IndexOfStringArray(methods, apiInfo[0]) == -1 {
				GoCFLog(apiInfo[0] + " is not currect")
				continue
			}

			path := parent + "/" + apiToPath(apiInfo[1])
			api := ScriptApi{
				Path:   path,
				Module: apiToPath(path),
				Method: apiInfo[0],
				File:   root + "/" + name,
			}

			ScriptApiMap = append(ScriptApiMap, api)
		}
	}

	return nil
}

// ******************** Api 初始化 End *************************
