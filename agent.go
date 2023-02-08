// 用于同步Master配置的脚本
package gocf

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryouaki/koa"
)

// 更新脚本文件接口参数结构定义
type ScriptParams struct {
	Files []ScriptFile `json:"files"`
}
type ScriptFile struct {
	Path   string `json:"path"`
	Script string `json:"script"`
}

// end

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
	agent.Post("/mapi/scripts", doSyncScripts)
}

func doCheck(ctx *koa.Context, next koa.Next) {
	ctx.SetBody(buildResp(false, "", ""))
}

/**
 * 增量同步脚本接口
 * {
 *		files: [
 *				{
 *						path: '',
 *						script: ''
 *				}
 *    ]
 * }
 */
func doSyncScripts(ctx *koa.Context, next koa.Next) {
	if ctx.Body == nil {
		ctx.Status = 400
		ctx.SetBody(buildResp(true, "Params Error", "参数为空"))
		return
	}

	data := ScriptParams{}
	err := json.Unmarshal(ctx.Body, &data)
	if err != nil {
		ctx.Status = 400
		ctx.SetBody(buildResp(true, "Params Error", "参数格式错误"))
		return
	}

	if len(data.Files) <= 0 {
		ctx.SetBody(buildResp(false, "", ""))
		return
	}

	scriptDevDir := strings.Replace(Root, "scripts", "scripts_dev", -1)
	os.RemoveAll(scriptDevDir)
	os.MkdirAll(scriptDevDir, 0750)

	// 拷贝文件到临时目录
	CopyTo(Root, scriptDevDir)

	// 做增量覆盖
	err = ReplaceScript(data, scriptDevDir)
	if err != nil {
		ctx.Status = 500
		ctx.SetBody(buildResp(true, "Server Error", "解析错误，请重试"))
		return
	}

	// 重置dev VM
	ClearApiMap(true)
	err = LoadApiScripts(scriptDevDir+"/api", true, "/api/dev")
	if err != nil {
		ctx.Status = 500
		ctx.SetBody(buildResp(true, "Server Error", "开发环境加载失败，请重试"))
		return
	}

	InitDevVM()
	err = InitApi(true)
	if err != nil {
		ctx.Status = 500
		ctx.SetBody(buildResp(true, "Server Error", "开发环境加载失败，请重试"))
		return
	}

	ctx.SetBody(buildResp(false, "", "文件同步成功"))
}

func ReplaceScript(files ScriptParams, distDir string) error {
	for _, f := range files.Files {
		filename := distDir + "/" + f.Path
		dir := filepath.Dir(filename)
		err := os.MkdirAll(dir, 0750)
		if err != nil && !os.IsExist(err) {
			return err
		}
		data := []byte(f.Script)
		err = os.WriteFile(filename, data, 0750)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildResp(err bool, msg string, data interface{}) []byte {
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
