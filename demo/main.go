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
	plugins := make([]*gocf.PluginCb, 0, 4)
	plugins = initConsole(plugins)
	gocf.RegistPlugin("console", plugins)

	gocf.InitGoCloudFunc()
}

func initConsole(plugins []*gocf.PluginCb) []*gocf.PluginCb {
	plugin := new(gocf.PluginCb)
	plugin.Name = "log"
	plugin.Fb = func(args []*gocf.JSValue, this *gocf.JSValue) (*gocf.JSValue, *gocf.JSValue) {
		goArgs := make([]any, 0, 4)
		for _, v := range args {
			val := v.ToString()
			goArgs = append(goArgs, val)
		}
		fmt.Println(goArgs...)
		return nil, nil
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
