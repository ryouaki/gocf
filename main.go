package main

import (
	"fmt"

	"github.com/ryouaki/gocf/core"
	"github.com/ryouaki/koa"
)

func init() {
	core.InitGoCloudFunc()
}

func main() {
	app := koa.New()

	app.Get("/", func(ctx *koa.Context, next koa.Next) {

	})

	err := app.Run(8080)
	if err != nil {
		fmt.Println("Server Failed:", err)
	}
}
