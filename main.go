package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("", r.Index)
}

func (r *Repository) Index(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Hello, World 👋!",
	})
}

func main() {
	// Read env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	r := Repository{DB: nil}

	// Initialize Fiber
	app := fiber.New()

	// Routes
	r.SetupRoutes(app)

	// Run server on port 8000

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	fmt.Println("\nServer running on port", port)

	app.Listen(host + ":" + port)
}
