package plugins

import (
	"github.com/ryouaki/gocf"
)

func initHttp() []*gocf.Plugin {
	plugins := make([]*gocf.Plugin, 0, 4)
	plugin := makePlugin("request", func(args []*gocf.JSValue, this *gocf.JSValue) *gocf.JSValue {
		// method := args[0]
		// if !method.IsString() {
		// 	return nil
		// }
		// uri := args[1]
		// if !uri.IsString() {
		// 	return nil
		// }

		// http.NewRequest(method.ToString(), uri.ToString())
		// v := make(map[string]string)
		// v["a"] = "aaaaa"
		// return gocf.MakeInvokeResult(this.Ctx, gocf.CB_SUCCESS, "aaaaaa")
		return nil
	})

	plugins = append(plugins, plugin)
	return plugins
}
