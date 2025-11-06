package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"washwise/config"
	"washwise/cron"
	"washwise/model"
	"washwise/server"
	"washwise/util"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := config.Load("config/config.yaml"); err != nil {
		fmt.Fprintf(os.Stderr, "加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	cfg := config.Get()

	// 使用配置初始化日志系统
	fmt.Println("正在初始化日志系统...")
	util.InitLogger(util.LogConfig{
		Level: cfg.Log.Level,
		Dir:   cfg.Log.Dir,
	})

	// 初始化数据库
	log.Info("初始化数据库...")
	if err := model.InitDB(cfg.Database.Path); err != nil {
		log.WithError(err).Fatal("初始化数据库失败")
	}
	log.Info("数据库初始化成功")

	// 初始化并启动定时任务
	taskManager := cron.InitTaskManager()
	taskManager.Start()

	// 初始化并启动HTTP服务器
	log.Info("初始化 HTTP 服务器...")
	srv := server.New(cfg)

	// 在goroutine中启动服务器
	go func() {
		if err := srv.Start(); err != nil {
			log.WithError(err).Fatal("HTTP 服务器启动失败")
		}
	}()

	// 等待终止信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 服务启动信息直接打印到标准输出
	fmt.Printf("  WashWise 服务已启动\n")
	fmt.Printf("  按 Ctrl+C 退出\n")
	<-sigChan

	// 服务关闭信息直接打印到标准输出
	fmt.Println("\n收到终止信号，正在关闭服务...")

	taskManager.Stop()
	if err := srv.Stop(); err != nil {
		fmt.Fprintf(os.Stderr, "关闭 HTTP 服务器失败: %v\n", err)
	}

	fmt.Println("服务已完全关闭")
}
