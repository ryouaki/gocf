/**
* 初始化加载api脚本
 */
package gocf

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type ScriptApi struct {
	Path   string // api地址
	Module string // 模块名
	Method string // api方法
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
var ScriptDevApiMap = make([]ScriptApi, 0, 4)

func ClearApiMap(isDev bool) {
	if isDev {
		ScriptDevApiMap = ScriptDevApiMap[0:0]
	} else {
		ScriptApiMap = ScriptApiMap[0:0]
	}
}

// 根据参数指定目录加载脚本
func LoadApiScripts(root string, isDev bool, parent string) error {
	GoCFLog(root)

	// 遍历脚本文件目录
	dir, err := ioutil.ReadDir(root)
	if err != nil {
		GoCFLog("Error", "Load Script file failed!")
		return err
	}

	for _, f := range dir {
		name := f.Name()
		if f.IsDir() {
			LoadApiScripts(root+"/"+name, isDev, parent+"/"+name)
		} else if !strings.HasSuffix(name, "js") {
			continue
		} else {
			apiInfo := strings.Split(name, ".")
			if IndexOfStringArray(methods, apiInfo[0]) == -1 {
				GoCFLog(apiInfo[0] + " is not currect")
				continue
			}

			path := parent + "/" + apiToPath(apiInfo[1])
			api := ScriptApi{
				Path:   path,
				Module: apiToPath(path),
				Method: apiInfo[0],
				File:   root + "/" + name,
			}

			if isDev {
				ScriptDevApiMap = append(ScriptApiMap, api)
			} else {
				ScriptApiMap = append(ScriptApiMap, api)
			}
		}
	}

	return nil
}

// 配置路径转访问路由路径
func apiToPath(path string) string {
	return strings.ReplaceAll(path, "/", "_")
}

func GetApiModule(method string, path string) (string, error) {
	var scripts []ScriptApi
	if strings.HasPrefix(path, "/api/dev") {
		scripts = ScriptDevApiMap
	} else {
		scripts = ScriptApiMap
	}
	for _, v := range scripts {
		if method == v.Method && path == v.Path {
			return v.Module, nil
		}
	}
	return "", fmt.Errorf("Api not found!")
}
