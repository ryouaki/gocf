// 用于同步Master配置的脚本
package gocf

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	GoCFLog("Din", ret)
}

func InitAgent(agent *koa.Application) {
	agent.Get("/mapi/check", doCheck)
	agent.Post("/mapi/scripts", doSyncScripts)
	agent.Post("/mapi/restart", doRestart)
}

func doCheck(ctx *koa.Context, next koa.Next) {
	ctx.SetBody(buildResp(false, "", ""))
}

func doRestart(ctx *koa.Context, next koa.Next) {
	go doResetVM()
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

	scriptTmp := strings.Replace(Root, "scripts", "scripts_tmp", -1)
	os.RemoveAll(scriptTmp)
	os.MkdirAll(scriptTmp, 0750)

	// 拷贝文件到临时目录
	CopyTo(Root, scriptTmp)

	// 做增量覆盖
	err = ReplaceScript(data, scriptTmp)
	if err != nil {
		ctx.Status = 500
		ctx.SetBody(buildResp(true, "Server Error", "解析错误，请重试"))
		return
	}

	os.RemoveAll(Root)
	os.MkdirAll(Root, 0750)

	CopyTo(scriptTmp, Root)
	os.RemoveAll(scriptTmp)

	ClearApiMap()

	go doResetVM()

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

func doResetVM() {
	GoCFLog("doResetVM", "Restart Ok")
	updateResetFlag(true)
	for len(vms) > 0 {
		vm := vms[0]
		if vm.IsFree { // 需要等待所有处理结束后是否内存
			vm.Ctx.Free()
			vm.VM.Free()
			vms = vms[1:]
		}
		time.Sleep(time.Duration(1) * time.Millisecond)
	}
	vms = vms[0:0]
	RunGoCF()
	updateResetFlag(false)
	GoCFLog("doResetVM", "Restart End")
}
