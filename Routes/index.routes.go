package Routes

import (
	"github.com/Data-Alchemist-ODS/ods-api/Controllers"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Post("/v1/api", Controllers.CreateTransaction)
	r.Get("/v1/api/transaction", Controllers.GetAllTransactions)
}
