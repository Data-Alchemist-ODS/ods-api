package routes

import (
	"github.com/Data-Alchemist-ODS/ods-api/controllers"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Post("/v1/api", controllers.CreateTransaction)
	r.Get("/v1/api/transaction", controllers.GetAllTransactions)
}
