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
func initConsole() []*gocf.Plugin {
	plugins := make([]*gocf.Plugin, 0, 1)
	plugin := makePlugin("log", func(args []*gocf.JSValue, this *gocf.JSValue) *gocf.JSValue {
		goArgs := make([]any, 1, 4)
		goArgs[0] = "[GoCF]:"
		for _, v := range args {
			val := v.ToString()
			goArgs = append(goArgs, val)
		}
		gocf.GoCFLog(goArgs...)
		return nil
	})

	plugins = append(plugins, plugin)

	return plugins
}
