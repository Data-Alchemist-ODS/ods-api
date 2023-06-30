package Routes 

import (
	"github.com/Data-Alchemist-ODS/ods-api/Controllers"

	"github.com/gofiber/fiber/v2"
)

func RouteInit (r *fiber.App) {
	r.Get("/", Controllers.GetALLTranscation)
}