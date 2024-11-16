package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Protected middleware
func Protect(c *fiber.Ctx) error {
	// 1) Get the token from the request header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed token"})
	}

	// 2) Parse the token
	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// 3) Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing claims"})
	}

	// 4) Set the user in the locals
	c.Locals("user", claims)
	return c.Next()
}

// Authorize checks if the user has the required role
func Authorize(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1) Get the user from the locals
		userClaims, ok := c.Locals("user").(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing claims"})
		}

		// 2) Check if the user has the required role
		userRole, ok := userClaims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error parsing role"})
		}

		for _, role := range roles {
			if role == userRole {
				return c.Next()
			}
		}

		// 3) Return an error if the user does not have the required role
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})

	}
}
