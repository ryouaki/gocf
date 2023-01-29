package controller

import (
	"fmt"

	"github.com/ryouaki/gocfm/render"
	"github.com/ryouaki/koa"
)

func InitController(app *koa.Application) {
	app.Get("", func(ctx *koa.Context, n koa.Next) {
		ctx.Status = 200
		tpl, err := render.RenderLogin(ctx)
		if err != nil {
			fmt.Println("111")
		}
		ctx.SetBody(tpl)
	})
}
