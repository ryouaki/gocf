package render

import (
	"fmt"

	"github.com/ryouaki/koa"
)

func RenderLogin(ctx *koa.Context) ([]byte, error) {
	tpl, err := buildTemplate("login.html", ctx, base+"template/login.html", base+"template/common/head.html")
	if err != nil {
		errString := fmt.Sprintf("System error %s", err)
		tempData := []byte(errString)
		return tempData, err
	}

	return tpl, nil
}
