package servicev2

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(r fiber.Router) {
	r.Get("/shops", GetShops)
	r.Get("/machines", GetMachines)
	r.Get("/machine/:machineId", GetMachine)
	r.Get("/machine/:machineId/like", Like)
	r.Get("/machine/:machineId/dislike", DisLike)
}
