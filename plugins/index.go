package plugins

import "github.com/ryouaki/gocf"

func InitPlugins() {
	// gocf.RegistPlugin("console", initConsole())
	// gocf.RegistPlugin("http", initHttp())
}

func makePlugin(name string, cbFunc gocf.JSGoFuncHandler) *gocf.Plugin {
	plugin := new(gocf.Plugin)
	plugin.Name = name
	plugin.Cb = cbFunc
	return plugin
}
