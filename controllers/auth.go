package controllers

import (
	"context"
	"os"
	"time"

	"github.com/JongSinister/WTFiber/config"
	"github.com/JongSinister/WTFiber/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const userCollection = "users"

// @desc	Register a new user
// @route	POST /api/v1/auth/register
// @access	Public
func Register(c *fiber.Ctx) error {

	// 1) Parse the request body into the User struct
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 2) Validate the user input
	if !user.ValidateEmail() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email address"})
	}

	// 3) check if the email is already registered
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := config.DB.Collection(userCollection).CountDocuments(ctx, bson.M{"email": user.Email})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking email"})
	}

	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already registered"})
	}

	// 4) Hash the user's password
	if err := user.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}

	user.CreatedAt = primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))

	// 5) Insert the user into the database
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := config.DB.Collection(userCollection).InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	// 6) Generate a JWT token for the user
	token, err := user.GenerateToken(os.Getenv("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating token"})
	}

	// 7) Return the response
	return SendCookie(c, fiber.StatusOK, token, user.ID)
}

// @desc	Login a user
// @route	POST /api/v1/auth/login
// @access	Public
func Login(c *fiber.Ctx) error {
	// 1) Parse the request body into the User struct
	loggedUser := new(models.User)

	if err := c.BodyParser(&loggedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// 2) Validate the email format
	if !loggedUser.ValidateEmail() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email address"})
	}

	// 3) Find the user by email
	targetUser := new(models.User)
	err := config.DB.Collection(userCollection).FindOne(context.Background(), bson.M{"email": loggedUser.Email}).Decode(targetUser)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	// 4) Check the password
	if !targetUser.CheckPassword(loggedUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
	}

	// 5) Generate a JWT token for the user
	token, err := targetUser.GenerateToken(os.Getenv("JWT_SECRET"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	// 6) Return the response
	return SendCookie(c, fiber.StatusOK, token, targetUser.ID)
}

// @desc	Get the current user
// @route	GET /api/v1/auth/me
// @access	Private
func Me(c *fiber.Ctx) error {
	// 1) Get the user's email from the JWT claims
	userClaims, ok := c.Locals("user").(jwt.MapClaims)["email"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching user data"})
	}

	// 2) Find the user by email
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := new(models.User)
	err := config.DB.Collection(userCollection).FindOne(ctx, bson.M{"email": userClaims}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching user data"})
	}

	// 3) Return the response
	return c.Status(fiber.StatusOK).JSON(user)
}

// @desc    Log user out / clear cookie
// @route   GET /api/v1/auth/logout
// @access  Private
func Logout(c *fiber.Ctx) error {

	// 1) Create a cookie object with an expired date to clear it
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true, // Ensure the cookie is HttpOnly
	}

	// 2) Send the cookie to the client
	c.Cookie(&cookie)

	// 3) Respond with a success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"date":    time.Now().Format(time.RFC3339),
	})
}

// Send Cookie function
func SendCookie(c *fiber.Ctx, statusCode int, token string, userID primitive.ObjectID) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	})

	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"token":   token,
		"userid":  userID,
	})
}
