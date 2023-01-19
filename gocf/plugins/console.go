package plugins

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L../ -lquickjs

#include "../quickjs-libc.h"
*/
import "C"

import (
	"github.com/ryouaki/gocf"
)

// 初始化console插件
func initConsole() []*gocf.PluginCb {
	plugins := make([]*gocf.PluginCb, 0, 1)
	plugin := new(gocf.PluginCb)
	plugin.Name = "log"
	plugin.Fb = func(args []*gocf.JSValue, this *gocf.JSValue) *gocf.JSValue {
		goArgs := make([]any, 0, 4)
		for _, v := range args {
			val := v.ToString()
			goArgs = append(goArgs, val)
		}
		return nil
	}

	plugins = append(plugins, plugin)

	return plugins
}
