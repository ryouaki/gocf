package gocfm

import "github.com/ryouaki/koa/log"

func init() {
	// 初始化日志系统
	log.New(&log.Config{
		Level:   log.LevelInfo, // 日志级别info
		Mode:    log.LogStd,    // 以文件形式保存日志
		MaxDays: 1,             // 最多保留一天
		LogPath: "./logs",      // 日志文件存储位置
	})
}

func InitLog() {
	log.Info("Log Start OK")
}
