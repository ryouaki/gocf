package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ryouaki/gocf"
)

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "../quickjs-libc.h"
*/
import "C"

func init() {
	gocf.RegistPlugin("console", initConsole())
	gocf.RegistPlugin("http", initHttp())

	gocf.InitGoCloudFunc()
}

func initConsole() []*gocf.PluginCb {
	plugins := make([]*gocf.PluginCb, 0, 4)
	plugin := new(gocf.PluginCb)
	plugin.Name = "log"
	plugin.Fb = func(args []*gocf.JSValue, this *gocf.JSValue) *gocf.JSValue {
		goArgs := make([]any, 0, 4)
		for _, v := range args {
			val := v.ToString()
			goArgs = append(goArgs, val)
		}
		fmt.Println(goArgs...)
		return nil
	}

	plugins = append(plugins, plugin)

	return plugins
}

func initHttp() []*gocf.PluginCb {
	plugins := make([]*gocf.PluginCb, 0, 4)
	plugin := new(gocf.PluginCb)
	plugin.Name = "request"
	plugin.Fb = func(args []*gocf.JSValue, this *gocf.JSValue) *gocf.JSValue {
		method := args[0]
		if !method.IsString() {
			return nil
		}
		uri := args[1]
		if !uri.IsString() {
			return nil
		}

		// http.NewRequest(method.ToString(), uri.ToString())
		return nil
	}

	plugins = append(plugins, plugin)
	return plugins
}

func main() {
	f, err := os.OpenFile("./main.js", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Read Script Failed:", err)
		f.Close()
		return
	}

	src, err1 := ioutil.ReadAll(f)
	if err1 != nil {
		fmt.Println("Read Script Failed:", err1)
	}

	gocf.RunAPI(string(src))
}
