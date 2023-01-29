package main

import (
	"fmt"

	"github.com/ryouaki/gocfm"
	"github.com/ryouaki/koa"
)

func main() {
	app := koa.New()

	gocfm.InitController(app)

	err := app.Run(8001) // 启动
	if err != nil {      // 是否发生错误
		fmt.Println(err)
	}

}
