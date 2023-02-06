package main

import (
	"fmt"
	"os"

	"github.com/ryouaki/gocfm"
	"github.com/ryouaki/gocfm/controller"
	"github.com/ryouaki/koa"
	"github.com/ryouaki/koa/session"
	"github.com/ryouaki/koa/static"
)

func main() {
	app := koa.New()

	if os.Getenv("MODE") != "DEV" {
		app.Use(static.Static("./static", "/static/"))
	}

	sessConf := session.SessionConf{
		MaxAge: 60 * 30,
	}
	app.Use(session.Session(sessConf, gocfm.SessionStore))

	gocfm.InitLog()

	app.Use(func(ctx *koa.Context, next koa.Next) {
		sess := ctx.GetData("session").(map[string]interface{})
		if sess["isLogin"] == nil && ctx.Path != "/api/dologin" {
			ctx.Status = 401
			ctx.SetBody([]byte("Please sign in first"))
		} else {
			ctx.SetHeader("Content-Type", "application/json")
			next()
		}
	})
	controller.Init(app)

	err := app.Run(8001) // 启动
	if err != nil {      // 是否发生错误
		fmt.Println(err)
	}

}
