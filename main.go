// main.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/Routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
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
	Routes.RouteInit(app)

	//handle unavailavble route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
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
