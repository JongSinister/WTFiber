package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// Hotel routes
	HotelRoutes(api.Group("/hotels"))
}