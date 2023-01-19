package plugins

import "github.com/ryouaki/gocf"

func InitPlugins() {
	gocf.RegistPlugin("console", initConsole())
	gocf.RegistPlugin("http", initHttp())
}

func makePlugin(name string, cbFunc gocf.JSGoFuncHandler) *gocf.PluginCb {
	plugin := new(gocf.PluginCb)
	plugin.Name = name
	plugin.Fb = cbFunc
	return plugin
}
