/**
* 初始化加载api脚本
 */
package gocf

import (
	"io/ioutil"
	"strings"
)

type ScriptApi struct {
	Path   string // api地址
	Method string // api方法
	Ver    string // api脚本版本
	File   string // 脚本文件地址
}

const (
	METHOD_GET     = "get"
	METHOD_POST    = "post"
	METHOD_OPTIONS = "options"
	METHOD_PUT     = "put"
	METHOD_DELETE  = "delete"
	METHOD_PATCH   = "patch"
	METHOD_HEAD    = "head"
)

var methods = []string{
	METHOD_GET,
	METHOD_POST,
	METHOD_OPTIONS,
	METHOD_PUT,
	METHOD_DELETE,
	METHOD_PATCH,
	METHOD_HEAD,
}

var ScriptApiMap = make([]ScriptApi, 0, 4)

// 根据参数指定目录加载脚本
func LoadApiScripts(root string) {
	GoCFLog(root)

	// 遍历脚本文件目录
	dir, err := ioutil.ReadDir(root)
	if err != nil {
		GoCFLog("Error", "Load Script file failed!")
		return
	}

	for _, f := range dir {
		name := f.Name()
		// 过滤掉非js结尾的文件和目录
		if f.IsDir() || !strings.HasSuffix(name, "js") {
			continue
		}
		apiInfo := strings.Split(name, ".")
		// 只有以default结尾的api才是发布的api
		if len(apiInfo) < 4 || apiInfo[3] != "default" {
			continue
		}
		if IndexOfStringArray(methods, apiInfo[0]) == -1 {
			continue
		}
		api := ScriptApi{
			Path:   apiToPath(apiInfo[1]),
			Method: apiInfo[0],
			Ver:    apiInfo[2],
			File:   root + "/" + name,
		}
		ScriptApiMap = append(ScriptApiMap, api)
	}
}

// 配置路径转访问路由路径
func apiToPath(path string) string {
	return strings.ReplaceAll(path, "_", "/")
}
