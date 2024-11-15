package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/JongSinister/WTFiber/config"
	"github.com/JongSinister/WTFiber/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	envPath := filepath.Join(cwd, "../config/.env")
	log.Printf("Loading .env file from: %s", envPath)

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	app := fiber.New()

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	// Connect to MongoDB
	config.ConnectDB()
	defer config.DisconnectDB()

	// Initialize the database
	config.InitDB()

	// set up routes
	routes.Setup(app)

	port := os.Getenv("PORT")
	log.Fatal(app.Listen("localhost:" + port))
}

//nodemon --watch '*.go' --exec 'go run main.go' --ext go
