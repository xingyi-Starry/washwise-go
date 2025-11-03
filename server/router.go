package server

import (
	servicev1 "washwise/server/service_v1"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RegisterServices(app *fiber.App) {
	app.Use(cors.New(cors.ConfigDefault))

	servicev1.RegisterRoutes(app.Group("/api/v1"))
}
