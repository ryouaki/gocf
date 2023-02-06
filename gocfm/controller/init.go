package controller

import (
	"github.com/ryouaki/koa"
)

func Init(app *koa.Application) {
	app.Post("/api/dologin", doLogin)
}
