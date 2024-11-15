package controllers

import (
	"context"
	"log"
	"time"

	"github.com/JongSinister/WTFiber/config"
	"github.com/JongSinister/WTFiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const hotelCollection = "hotels"

// @desc    Get all hotels
// @route   GET /api/v1/hotels/
// @access  Private
func GetHotels(c *fiber.Ctx) error {
	// 1) Prepare the query to fetch all hotels
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2) Fetch all hotels from the database
	opts := options.Find()
	cursor, err := config.DB.Collection(hotelCollection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error fetching hotels"})
	}
	defer cursor.Close(ctx)

	var hotels []models.Hotel
	if err := cursor.All(ctx, &hotels); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Error fetching hotels",
			"msg":   err,
		})
	}

	log.Printf("Fetched hotels: %+v", hotels)

	if len(hotels) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No hotels found"})
	}

	return c.JSON(hotels)

}
