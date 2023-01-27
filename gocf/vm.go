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
	"strconv"
	"sync"
	"time"
)

// 启动引擎数量，默认1
var Nums = 1

// 脚本加载目录，默认./
var Root = "./"

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
var vms = make([]JSVM, 0, 4)

// Go插件缓存池
var pluginMap = make(map[string][]*Plugin)

func InitVM(nums int) {
	// 如果为0则直接退出
	if nums < 1 {
		return
	}
	// 根据启动参数-n初始化引擎数量
	for i := 0; i < nums; i++ {
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

		vms = append(vms, JSVM{
			VM:     rt,
			Ctx:    ctx,
			IsFree: true,
		})
	}
	GoCFLog("initialized " + strconv.Itoa(nums) + " JSVM")
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

	var wg sync.WaitGroup

	wg.Add(1)

	var vm *JSVM = nil

	go func() {
		for {
			for _, v := range vms {
				if v.IsFree {
					v.IsFree = false
					vm = &v
					wg.Done()
				}
			}
			// 不能一直占用cpu，需要休眠1ms
			time.Sleep(time.Duration(1) * time.Millisecond)
		}
	}()
	go func() {
		select {
		case <-ctx.Done():
			vm = nil
			wg.Done()
		}
	}()

	// 等待2个协程其中之一返回
	wg.Wait()

	return vm
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
