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
	//Transaction Routes
	r.Get("/v1/api/transaction", transactionController.GetAllTransactions)
	r.Get("/v1/api/transaction/:id", transactionController.GetOneTransaction)

	r.Get("/v1/api/data", transactionController.GetAllStoredDatas)
	r.Get("/v1/api/data/:id", transactionController.GetOneStoredData)

	//POST ROUTES
	//Transaction Routes
	r.Post("/v1/api/transaction", transactionController.CreateNewTransaction)
	r.Post("/v1/api/connect/tidb", databaseController.ConnectToTiDB)

	//DELETE ROUTES
	//Transaction Routes
	r.Delete("/v1/api/transaction/:id", transactionController.DeleteTransaction)
}
