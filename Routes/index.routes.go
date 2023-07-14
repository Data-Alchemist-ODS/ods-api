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

	//GET ROUTES
	r.Get("/v1/api/transaction", transactionController.GetAllTransactions)
	r.Get("/v1/api/transaction/:id", transactionController.GetOneTransaction)
	r.Get("/v1/api/data", transactionController.GetAllStoredDatas)

	//POST ROUTES
	r.Post("/v1/api/transaction", transactionController.CreateNewTransaction)
	r.Post("/v1/api/connect/tidb", databaseController.ConnectToTiDB)
}
