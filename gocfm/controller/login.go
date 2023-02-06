package controller

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/ryouaki/gocfm"
	"github.com/ryouaki/koa"
)

func doLogin(ctx *koa.Context, next koa.Next) {
	body := ctx.Body
	if body == nil {
		ctx.Status = 404
		ctx.SetBody([]byte("参数错误，登录失败"))
		return
	}

	jsonData := make(map[string]interface{})
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		ctx.Status = 404
		ctx.SetBody([]byte("参数错误，登录失败"))
		return
	}

	db := gocfm.GetDB()
	pwd := gocfm.InterfaceToString(jsonData["Password"])
	rows, err := db.Query(fmt.Sprintf(`select UserName,Password from Users where UserName='%s' and Password='%s'`, jsonData["UserName"], md5.Sum([]byte(pwd))))
	defer rows.Close()
	if err != nil {
		ctx.Status = 500
		ctx.SetBody([]byte("系统错误，请稍后再试"))
		return
	}

	if rows.Next() {
		sess := ctx.GetData("session").(map[string]interface{})
		sess["isLogin"] = true
		ctx.SetData("session", sess)
		ctx.SetBody([]byte(`{"error": true, "data": "", "msg": "成功"}`))
		return
	}

	ctx.Status = 404
	ctx.SetBody([]byte("参数错误，登录失败"))
}
