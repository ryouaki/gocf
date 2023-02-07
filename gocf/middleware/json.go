package middleware

import "github.com/ryouaki/koa"

func initJson(ctx *koa.Context, next koa.Next) {
	next()
	ctx.SetHeader("Content-Type", "application/json")
}
