package routes

import (
	"github.com/Data-Alchemist-ODS/ods-api/controllers"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	transactionController := controllers.NewTransactionController()
	databaseController := controllers.NewDatabaseController()
	userController := controllers.NewUserController()

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to On-demand Data Sharding 👋",
		})
	})

	//GET ROUTES
	//User Routes
	r.Get("/v1/api/user", userController.GetAllUser)
	r.Get("/v1/api/user/:id", userController.GetOneUser)

	//Transaction Routes
	r.Get("/v1/api/transaction", transactionController.GetAllTransactions)
	r.Get("/v1/api/transaction/:id", transactionController.GetOneTransaction)

	r.Get("/v1/api/data", transactionController.GetAllStoredDatas)
	r.Get("/v1/api/data/:id", transactionController.GetOneStoredData)

	//POST ROUTES
	//User Routes
	r.Post("/v1/api/user", userController.CreateUser)

	//Transaction Routes
	r.Post("/v1/api/transaction", transactionController.CreateNewTransaction)
	r.Post("/v1/api/connect/tidb", databaseController.ConnectToTiDB)

	//PUT ROUTES
	//User Routes
	r.Put("/v1/api/user/update/:id", userController.UpdateUser)
	
	//DELETE ROUTES
	//User Routes
	r.Delete("/v1/api/user/:id", userController.DeleteOneUser)

	//Transaction Routes
	r.Delete("/v1/api/transaction/:id", transactionController.DeleteTransaction)
}
