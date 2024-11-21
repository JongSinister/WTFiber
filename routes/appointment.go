package routes

import (
	"github.com/JongSinister/WTFiber/controllers"
	"github.com/JongSinister/WTFiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func AppointmentRoutes(router fiber.Router) {
	router.Get("/", middleware.Protect, middleware.Authorize("admin"), controllers.GetAppointments)
	router.Get("/:id", middleware.Protect, middleware.Authorize("admin"), controllers.GetAppointment)
	router.Put("/:id", middleware.Protect, middleware.Authorize("admin"), controllers.UpdateAppointment)
	router.Delete("/:id", middleware.Protect, middleware.Authorize("admin"), controllers.DeleteAppointment)
}
