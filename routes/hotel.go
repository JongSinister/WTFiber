package routes

import (
	"github.com/JongSinister/WTFiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func HotelRoutes(router fiber.Router) {
	router.Get("/", controllers.GetHotels) // Relative to the base path "/api/v1/hotels"
	// router.Get("/:id", controllers.GetHotel)
	// router.Post("/", controllers.CreateHotel)
	// router.Put("/:id", controllers.UpdateHotel)
	// router.Delete("/:id", controllers.DeleteHotel)
}
