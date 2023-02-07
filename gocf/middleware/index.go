package middleware

import "github.com/ryouaki/koa"

func Init(app *koa.Application) {
	app.Use(initJson)
}
