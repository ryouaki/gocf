package controller

import "github.com/ryouaki/koa"

/**
* ip 客户端ip地址
* last_update 当前时间
* cpu cpu使用率
* mem 内存使用率
* status 默认1 活跃，大于5分钟即失活
* listen 是否使用中
 */
func registClient(ctx *koa.Context, next koa.Next) {

}
