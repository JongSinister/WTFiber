package routes

import (
	"github.com/JongSinister/WTFiber/controllers"
	"github.com/JongSinister/WTFiber/middleware"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router fiber.Router) {
	router.Post("/register", controllers.Register)
	router.Post("/login", controllers.Login)
	router.Get("/me", middleware.Protect, controllers.Me)
	router.Get("/logout", controllers.Logout)
}
