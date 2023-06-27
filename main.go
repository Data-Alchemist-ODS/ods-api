package main

import (
	"fmt"
	"log"
	// "net/http"
	"os"

	// "github.com/Data-Alchemist-ODS/ods-api/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	// "gorm.io/gorm"	

	"github.com/Data-Alchemist-ODS/ods-api/database"
)

// type Repository struct {
// 	DB *gorm.DB
// }

// func (r *Repository) SetupRoutes(app *fiber.App) {
//	api := app.Group("/api")

// 	api.Get("", r.Index)
// }

// func (r *Repository) Index(c *fiber.Ctx) error {
// 	return c.Status(http.StatusOK).JSON(fiber.Map{
// 		"message": "Hello, World 👋!",
// 	})
// }

func main() {
	// Read env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	err = database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Fiber
	app := fiber.New()

	app.Get("/:name", func(ctx *fiber.Ctx) error{
		return ctx.SendString("Hello, " + ctx.Params("name") + "!")
	})

	// Routes
	// r.SetupRoutes(app)

	// Run server on port 8000

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	fmt.Println("\nServer running on port", port)

	app.Listen(host + ":" + port)
}
