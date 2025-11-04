package servicev1

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	app.Get("/getLaundryMachines", GetLaundryMachines)
	app.Get("/getMachineDetail", GetMachineDetail)
}
