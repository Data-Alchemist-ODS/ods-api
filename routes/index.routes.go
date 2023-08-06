package routes

import (
	"github.com/Data-Alchemist-ODS/ods-api/controllers"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	transactionController := controllers.NewTransactionController()
	databaseController := controllers.NewDatabaseController()
	userController := controllers.NewUserController()
	queryController := controllers.NewNaturalQueryController()
	notesController := controllers.NewNotesController()

	r.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to On-demand Data Sharding API Public Endpoint👋",
			"status" : fiber.StatusOK,
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

	//Notes Routes
	r.Get("/v1/api/gen-notes", notesController.GetAllNotes)

	//POST ROUTES
	//User Routes
	r.Post("/v1/api/user/register", userController.RegisterUser)
	r.Post("/v1/api/user/login", userController.LoginUser)

	//Transaction Routes
	r.Post("/v1/api/transaction", transactionController.CreateNewTransaction)
	r.Post("/v1/api/connect/tidb", databaseController.ConnectToTiDB)

	//Natural Query Routes
	r.Post("/v1/api/query", queryController.CreateNaturalQuery)

	//Notes Routes
	r.Post("/v1/api/gen-notes/create", notesController.CreateNewNotes)

	//PUT ROUTES
	//User Routes
	r.Put("/v1/api/user/update/:id", userController.UpdateUser)

	//DELETE ROUTES
	//User Routes
	r.Delete("/v1/api/user/:id", userController.DeleteOneUser)

	//Transaction Routes
	r.Delete("/v1/api/transaction/:id", transactionController.DeleteTransaction)

	//Notes Routes
	r.Delete("/v1/api/gen-notes/delete/:id", notesController.DeleteNotesHistory)
}