package controllers

import (
	"context"
	"math/rand"
	"time"

	"github.com/JongSinister/WTFiber/config"
	"github.com/JongSinister/WTFiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const appointmentCollection = "appointments"

// @desc Get all appointments
// @route GET /api/v1/appointments
// @access Private
func GetAppointments(c *fiber.Ctx) error {
	// 1) Fetch appointments from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find()
	cursor, err := config.DB.Collection(appointmentCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error fetching appointments"})
	}
	defer cursor.Close(ctx)

	// 2) Decode results
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

	// 3) Return appointments
	return c.JSON(appointments)
}

// @desc Get a single appointment
// @route GET /api/v1/appointments/:id
// @access Private
func GetAppointment(c *fiber.Ctx) error {
	// 1) Get the ID from the URL and convert it to an ObjectID
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID Format"})
	}

	// 2) Fetch the appointment by ID from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3) fetch appointment from database
	appointment := models.Appointment{}
	err = config.DB.Collection(appointmentCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Appointment not found"})
	}

	// 4) Return appointment
	return c.JSON(appointment)
}

// @desc Add appointment
// @route POST /api/v1/hotel/:hotelId/appointments
// @access Private
func AddAppointment(c *fiber.Ctx) error {
	// 1) Get the hotel ID from the URL and convert it to an ObjectID
	hotelID := c.Params("hotelId")

	objectHotelID, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Hotel ID Format"})
	}

	// 2) Parse the request body into the Appointment struct
	appointment := new(models.Appointment)
	if err := c.BodyParser(appointment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 3) Set the hotel ID and CreatedAt fields
	appointment.Hotel = objectHotelID
	appointment.CreatedAt = primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))
	appointment.WifiPassword = generateRandomPassword()

	// 4) Insert the appointment into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := config.DB.Collection(appointmentCollection).InsertOne(ctx, appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create appointment"})
	}

	// 5) Return the response
	return c.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"message":     "Appointment created successfully",
			"appointment": res,
		},
	)
}

// Create Random Wifi Password
func generateRandomPassword() string {
	// Random length from 6 to 8
	length := rand.Intn(3) + 6 // Random number between 6 and 8
	allStr := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIdx := rand.Intn(len(allStr))
		password[i] = allStr[randomIdx]
	}

	return string(password)
}
