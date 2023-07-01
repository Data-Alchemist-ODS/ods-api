package Routes

import (
	"github.com/Data-Alchemist-ODS/ods-api/Controllers"
	"github.com/gofiber/fiber/v2"
)

func RouteInit(r *fiber.App) {
	r.Post("/", Controllers.CreateTransaction)
	r.Get("/transaction", Controllers.GetAllTransactions)
}
