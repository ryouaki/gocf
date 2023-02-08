// 用于同步Master配置的脚本
package gocf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/ryouaki/koa"
)

func DinMaster() {
	ret, err := http.NewRequest("GET", MasterHost+"/api/din", nil)
	if err != nil {
		GoCFLog("Din Master failed, please restart")
		return
	}
	GoCFLog(ret)
}

func InitAgent(agent *koa.Application) {

	agent.Use(func(ctx *koa.Context, next koa.Next) {
		next()
		ctx.SetHeader("Content-Type", "application/json")
	})
	agent.Get("/mapi/check", doCheck)
	agent.Post("/mapi/reset", doReset)
	agent.Post("/mapi/update", doUpdate)
}

func doCheck(ctx *koa.Context, next koa.Next) {
	ctx.SetBody(makeRet(true, "", ""))
}

func doReset(ctx *koa.Context, next koa.Next) {

}

/*
* === api ===
* method:string 方法名
* script:string 脚本代码
* api:string api地址
* === module ===
* name: string 模块名
* script: string 脚本代码
 */
func doUpdate(ctx *koa.Context, next koa.Next) {
	body := ctx.Body
	if body == nil {
		ctx.Status = 400
		ctx.SetBody(makeRet(true, "Params Error", ""))
		return
	}

	data := make(map[string]interface{})
	err := json.Unmarshal(body, &data)
	if err != nil {
		ctx.Status = 400
		ctx.SetBody(makeRet(true, "Params Error", ""))
		return
	}

	var isApi = true
	var method = ""
	var script = ""
	var name = ""
	var api = ""
	for k, v := range data {
		if k == "method" {
			method = InterfaceToString(v)
			method = strings.TrimSpace(method)
		} else if k == "script" {
			script = InterfaceToString(v)
			script = strings.TrimSpace(script)
		} else if k == "name" {
			isApi = false
			name = InterfaceToString(v)
		} else if k == "api" {
			api = InterfaceToString(v)
		} else {
			ctx.Status = 400
			ctx.SetBody(makeRet(true, "Params Error", k+" is not be need"))
			return
		}
	}

	if isApi && (method == "" || script == "" || api == "") {
		ctx.Status = 400
		ctx.SetBody(makeRet(true, "Params Error", ""))
		return
	} else if !isApi && (script == "" || name == "") {
		ctx.Status = 400
		ctx.SetBody(makeRet(true, "Params Error", ""))
		return
	}

	dir := mkDir(isApi)

	fileName := ""
	if isApi {
		fileName = fmt.Sprintf("%s.%s.js", method, strings.ReplaceAll(api, "/", "_"))
	} else {
		fileName = fmt.Sprintf("%s.js", strings.ReplaceAll(name, "/", "_"))
	}

	/* 校验语法*/
	rt := NewRuntime()
	c := rt.NewContext()

	v, e := c.Eval(script, "test", 1<<0|1<<5)
	defer func() {
		v.Free()
		e.Free()
		c.Free()
		rt.Free()
	}()
	if c.GetException() != nil {
		r := c.GetException()
		ctx.Status = 400
		ctx.SetBody(makeRet(true, "Script Error", r.ToString()))
		return
	}

	fmt.Println(c.FindModule("test"))
	/* 校验语法*/

	err = ioutil.WriteFile(dir+"/"+fileName, []byte(script), 0666)
	if err != nil {
		ctx.Status = 500
		ctx.SetBody(makeRet(true, "System Error", fileName+" Update failed"))
		return
	}

	ctx.SetBody(makeRet(false, "", "ok"))
}

func mkDir(isApi bool) string {
	name := ""
	if isApi {
		name = "/apis"
	} else {
		name = "/modules"
	}
	err := os.Mkdir(Root+name, 0750)
	if err != nil && !os.IsExist(err) && !os.IsNotExist(err) {
		GoCFLog(err)
	}

	return Root + name
}

func checkIsExist(name string) bool {
	_, err := os.Stat(Root + "/" + name)
	if os.IsExist(err) {
		return true
	}
	return false
}

func makeRet(err bool, msg string, data interface{}) []byte {
	e := ""
	if err {
		e = `"error":true`
	} else {
		e = `"error":false`
	}

	m := ""
	if msg != "" {
		m = `"msg":"` + msg + `"`
	} else {
		m = `"msg":""`
	}

	d := ""
	if data != nil {
		d = `"data":` + InterfaceToString(data)
	} else {
		d = `"data":null`
	}

	return []byte("{" + e + "," + m + "," + d + "}")
}
