package controllers

import (
	"context"
	"time"

	"github.com/JongSinister/WTFiber/config"
	"github.com/JongSinister/WTFiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const appointmentCollection = "appointments"

// @desc Get all appointments
// @route GET /api/v1/appointments
// @access Private
func GetAppointments(c *fiber.Ctx) error {
	// 1) Retrieve user details from locals
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing claims"})
	}

	userRole, ok := userClaims["role"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing role"})
	}

	// 2) Check if the user has the required role
	if userRole == "user" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Access denied"})
	}

	// 3) Fetch appointments from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find()
	cursor, err := config.DB.Collection(appointmentCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error fetching appointments"})
	}
	defer cursor.Close(ctx)

	// 4) Decode results
	var appointments []models.Appointment
	if err := cursor.All(ctx, &appointments); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Error fetching appointments",
			"msg":   err,
		})
	}

	if len(appointments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No appointments found"})
	}

	// 5) Return appointments
	return c.JSON(appointments)
}
