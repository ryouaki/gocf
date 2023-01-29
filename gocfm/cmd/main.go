package main

import (
	"fmt"
	"os"

	"github.com/ryouaki/gocfm"
	"github.com/ryouaki/gocfm/controller"
	"github.com/ryouaki/gocfm/render"
	"github.com/ryouaki/koa"
	"github.com/ryouaki/koa/session"
	"github.com/ryouaki/koa/static"
)

func main() {
	app := koa.New()

	if os.Getenv("MODE") == "DEV" {
		app.Use(static.Static("../static", "/static/"))
	} else {
		app.Use(static.Static("./static", "/static/"))
	}

	sessConf := session.SessionConf{
		MaxAge: 60 * 30,
	}
	app.Use(session.Session(sessConf, gocfm.SessionStore))

	app.Use(func(ctx *koa.Context, next koa.Next) {
		isLogin := ctx.GetCookie("isLogin")
		if isLogin == nil {
			ctx.Status = 401
			tpl, err := render.RenderLogin(ctx)
			if err != nil {
				fmt.Println("111", err)
			}
			ctx.SetBody(tpl)
		} else {
			next()
		}
	})
	controller.InitController(app)

	err := app.Run(8001) // 启动
	if err != nil {      // 是否发生错误
		fmt.Println(err)
	}

}
