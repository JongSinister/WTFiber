package routes

import (
	"github.com/JongSinister/WTFiber/controllers"
	"github.com/JongSinister/WTFiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func HotelRoutes(router fiber.Router) {
	router.Get("/", controllers.GetHotels)
	router.Get("/:id", controllers.GetHotel)
	router.Post("/", middleware.Protect, middleware.Authorize("admin"), controllers.CreateHotel)
	router.Put("/:id", middleware.Protect, middleware.Authorize("admin"), controllers.UpdateHotel)
	router.Delete("/:id", controllers.DeleteHotel)

	// Create a appointment for a hotel
	router.Post("/:hotelId/appointments", middleware.Protect, controllers.AddAppointment)

}
