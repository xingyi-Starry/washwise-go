package server

import (
	servicev1 "washwise/server/service_v1"
	servicev2 "washwise/server/service_v2"

	_ "washwise/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func RegisterServices(app *fiber.App) {
	// middleware
	app.Use(cors.New(cors.ConfigDefault))

	// swagger
	app.Get("/docs/*", fiberSwagger.WrapHandler)

	// routes
	servicev1.RegisterRoutes(app.Group("/api/v1"))
	servicev2.RegisterRoutes(app.Group("/api/v2"))
}
