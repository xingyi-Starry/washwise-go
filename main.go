package main

import (
	"os"
	"os/signal"
	"syscall"
	"washwise/config"
	"washwise/cron"
	"washwise/model"

	log "github.com/sirupsen/logrus"
)

func main() {
	// 设置日志格式
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	// 加载配置文件
	log.Info("加载配置文件...")
	if err := config.Load("config/config.yaml"); err != nil {
		log.WithError(err).Fatal("加载配置文件失败")
	}

	cfg := config.Get()

	// 初始化数据库
	log.Info("初始化数据库...")
	if err := model.InitDB(cfg.Database.Path); err != nil {
		log.WithError(err).Fatal("初始化数据库失败")
	}
	log.Info("数据库初始化成功")

	// 初始化并启动定时任务
	taskManager := cron.InitTaskManager()
	taskManager.Start()

	// 等待终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Info("WashWise 服务已启动，按 Ctrl+C 退出")
	<-sigChan

	// 优雅关闭
	log.Info("收到终止信号，正在关闭...")
	taskManager.Stop()
	log.Info("服务已关闭")
}
