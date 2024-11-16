package routes

import (
	"github.com/JongSinister/WTFiber/controllers"
	"github.com/JongSinister/WTFiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func AppointmentRoutes(router fiber.Router) {
	router.Get("/", middleware.Protect, middleware.Authorize("admin", "user"), controllers.GetAppointments)
	// app.Get("/:id", controllers.GetAppointment)
	// app.Post("/", controllers.CreateAppointment)
	// app.Put("/:id", controllers.UpdateAppointment)
	// app.Delete("/:id", controllers.DeleteAppointment)
}
