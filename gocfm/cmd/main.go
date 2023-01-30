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

	app.Use(func(ctx *koa.Context, next koa.Next) {

	})
	controller.InitController(app)

	err := app.Run(8001) // 启动
	if err != nil {      // 是否发生错误
		fmt.Println(err)
	}

}
