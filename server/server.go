package server

import (
	"fmt"
	"time"
	"washwise/config"
	"washwise/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
)

const (
	readTimeout  = 30 * time.Second
	writeTimeout = 30 * time.Second
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

// New 创建新的服务器实例
func New(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	})

	// 添加中间件
	app.Use(recover.New())
	app.Use(util.FiberLogger())

	// 注册路由
	RegisterServices(app)

	return &Server{
		app: app,
		cfg: cfg,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	log.Infof("HTTP 服务器启动于 %s", addr)
	return s.app.Listen(addr)
}

// Stop 优雅关闭服务器
func (s *Server) Stop() error {
	log.Info("正在关闭 HTTP 服务器...")
	return s.app.Shutdown()
}
