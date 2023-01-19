package plugins

import "github.com/ryouaki/gocf"

func initHttp() []*gocf.PluginCb {
	plugins := make([]*gocf.PluginCb, 0, 4)
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
		v := make(map[string]string)
		v["a"] = "aaaaa"
		return gocf.MakeInvokeResult(this.Ctx, gocf.CB_SUCCESS, "aaaaaa")
	})

	plugins = append(plugins, plugin)
	return plugins
}
