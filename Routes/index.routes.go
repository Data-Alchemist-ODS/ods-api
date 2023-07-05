package routes

import (
	"github.com/Data-Alchemist-ODS/ods-api/controllers"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	transactionController := controllers.NewTransactionController()
	databaseController := controllers.NewDatabaseController()

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to On-demand Data Sharding 👋",
		})
	})

	r.Post("/v1/api", transactionController.CreateTransaction)
	r.Get("/v1/api/transaction", transactionController.GetAllTransactions)

	r.Post("/v1/api/connect/tidb", databaseController.ConnectToTiDB)
}
