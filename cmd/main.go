package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "../quickjs-libc.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"time"

	"github.com/ryouaki/gocf"
	"github.com/ryouaki/gocf/middleware"
	"github.com/ryouaki/gocf/plugins"
	"github.com/ryouaki/koa"
)

func init() {
	// 初始化系统插件
	plugins.InitPlugins()
	gocf.DinMaster()
}

func main() {
	gocf.RunGoCF()

	app := koa.New()

	middleware.Init(app)

	gocf.InitAgent(app)

	app.Use(func(ctx *koa.Context, next koa.Next) {
		// // 获取api对应的js模块
		moduleName, err := gocf.GetApiModule(ctx.Method, ctx.Path)
		if err != nil {
			ctx.Status = 404
			ctx.SetBody([]byte(err.Error()))
			return
		}

		rt := gocf.GetVM(time.Duration(1))

		if rt == nil {
			ctx.Status = 500
			ctx.SetBody([]byte("VM is busy now, Please retry again."))
			return
		}

		defer gocf.ReleaseVM(rt)

		var ret interface{} = nil
		var wg sync.WaitGroup
		wg.Add(1)

		// 实例化成功调用返回
		rejectCb := gocf.NewJSGoFunc(rt.Ctx, func(args []*gocf.JSValue, this *gocf.JSValue) *gocf.JSValue {
			ctx.Status = 400
			for _, v := range args {
				data := gocf.InterfaceToString(v)
				ret = data
				wg.Done()
				return nil
			}
			wg.Done()
			return nil
		})

		// 实例化成功调用返回
		resolveCb := gocf.NewJSGoFunc(rt.Ctx, func(args []*gocf.JSValue, this *gocf.JSValue) *gocf.JSValue {
			val := args[0]
			data := make(map[string]interface{})
			if val.GetProperty("error") != nil {
				data["error"] = val.GetProperty("error").ToString()
			}
			if val.GetProperty("data") != nil {
				data["data"] = val.GetProperty("data").ToString()
			}
			ret = data

			wg.Done()
			return nil
		})

		rt.Ctx.ExportFunc("resolve", resolveCb)
		rt.Ctx.ExportFunc("reject", rejectCb)

		method := ctx.Method
		query := gocf.InterfaceToString(ctx.Query)
		params := gocf.InterfaceToString(ctx.Params)
		body := string(ctx.Body)
		if body == "" {
			body = "\"\""
		}
		headers := gocf.InterfaceToString(ctx.Req.Header)

		exec := fmt.Sprintf("import exec from \"%s\";exec(\"%s\", %s, %s, %s, %s).then(resolve).catch(reject);", moduleName, method, query, params, body, headers)
		// exec := fmt.Sprintf("import exec from \"%s\";exec().then(resolve)", moduleName)
		fmt.Println(exec)
		wfb := rt.Ctx.Eval(exec, "<code>", 1)
		e := rt.Ctx.GetException()

		if e != nil && len(e.ToString()) > 0 {
			ctx.Status = 500
			ctx.SetBody([]byte(e.ToString()))
		} else {
			gocf.WaitForLoop(rt)
			wg.Wait()
			ctx.SetHeader("Content-Type", "application/json")
			data := ret
			ctx.SetBody([]byte(gocf.InterfaceToString(data)))
		}
		wfb.Free()
		e.Free()
	})

	err := app.Run(8000) // 启动
	if err != nil {      // 是否发生错误
		gocf.GoCFLog("Koa.Run", err)
	}
}
