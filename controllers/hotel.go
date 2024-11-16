package controllers

import (
	"context"
	"time"

	"github.com/JongSinister/WTFiber/config"
	"github.com/JongSinister/WTFiber/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const hotelCollection = "hotels"

// @desc    Get all hotels
// @route   GET /api/v1/hotels/
// @access  Public
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

	if len(hotels) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No hotels found"})
	}

	return c.JSON(hotels)

}

// @desc    Get a hotel by ID
// @route   GET /api/v1/hotels/:id
// @access  Public
func GetHotel(c *fiber.Ctx) error {

	// 1) Get the ID from the URL and convert it to an ObjectID
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID Format"})
	}

	// 2) Prepare the query to fetch the hotel by ID
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 3) Fetch the hotel by ID from the database
	hotel := models.Hotel{}

	err = config.DB.Collection(hotelCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&hotel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Hotel not found"})
	}

	return c.JSON(hotel)
}

// @desc    Create a new hotel
// @route   POST /api/v1/hotels/
// @access  Public
func CreateHotel(c *fiber.Ctx) error {
	// 1) Get user and check permission(do later)

	// 2) Parse the request body into a Hotel struct
	hotel := new(models.Hotel)
	if err := c.BodyParser(hotel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 3) Insert the hotel into the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := config.DB.Collection(hotelCollection).InsertOne(ctx, hotel)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create hotel"})
	}

	// 4) Return the response
	return c.Status(fiber.StatusCreated).JSON(res)
}

// @desc    Update a hotel by ID
// @route   PUT /api/v1/hotels/:id
// @access  Public
func UpdateHotel(c *fiber.Ctx) error {
	// 1) Get user and check permission(do later)

	// 2) Get the ID from the URL and convert it to an ObjectID
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID Format"})
	}

	// 3) Fetch the existing hotel document
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	existingHotel := new(models.Hotel)
	err = config.DB.Collection(hotelCollection).FindOne(ctx, bson.M{"_id": objectID}).Decode(existingHotel)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Hotel not found"})
	}

	// 4) Parse the request body into a map for partial updates
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 5) Prepare the update document
	update := bson.M{
		"$set": updates,
	}

	// 6) Update the hotel document with specified fields
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	updatedHotel := new(models.Hotel)
	err = config.DB.Collection(hotelCollection).FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(updatedHotel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update hotel"})
	}

	// 7) Return the updated hotel document
	return c.JSON(updatedHotel)
}

// @desc    Delete a hotel by ID
// @route   DELETE /api/v1/hotels/:id
// @access  Public
func DeleteHotel(c *fiber.Ctx) error {
	// 1) Get user and check permission(do later)

	// 2) Get the ID from the URL and convert it to an ObjectID
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID Format"})
	}

	// 3) Delete the hotel from database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := config.DB.Collection(hotelCollection).DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete hotel"})
	}

	if res.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "error"})
	}

	// 4) Return the response
	return c.JSON(fiber.Map{"message": "Hotel deleted successfully"})
}
