// 用于同步Master配置的脚本
package gocf

import (
	"net/http"

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
	sysInfo := make(map[string]interface{})

	sysInfo["ok"] = true

	data := InterfaceToString(sysInfo)

	ctx.SetBody([]byte(data))
}

func doReset(ctx *koa.Context, next koa.Next) {

}

func doUpdate(ctx *koa.Context, next koa.Next) {

}
