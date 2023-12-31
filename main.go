// main.go

package main

import (
	//default modules
	"fmt"
	"log"
	"os"

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//third-party modules
	"github.com/joho/godotenv"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/routes"
)

func main() {
	// Read env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	database.ConnectDB()

	defer database.DisconnectDB()

	// Initialize Fiber
	app := fiber.New()

	// Routes
	routes.RouteInit(app)

	//handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Route not found",
			"status":  fiber.StatusNotFound,
		})
	})

	// Run server on specified host and port
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	fmt.Println("\nServer running on", host+":"+port)

	err = app.Listen(host + ":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
