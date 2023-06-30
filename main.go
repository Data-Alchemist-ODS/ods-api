package main

import (
	"fmt"
	"log"
	// "net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/Routes"
)

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

	defer database.DisconnectDB()

	// Initialize Fiber
	app := fiber.New()

	app.Get("/:name", func(ctx *fiber.Ctx) error{
		return ctx.SendString("Hello, " + ctx.Params("name") + "!")
	})

	// Routes
	Routes.RouteInit(app)

	// Run server on port 8000

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	fmt.Println("\nServer running on port", port)

	app.Listen(host + ":" + port)
}
