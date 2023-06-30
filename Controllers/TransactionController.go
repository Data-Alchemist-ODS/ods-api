package Controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/Models/Entity"
)

func GetALLTranscation (ctx *fiber.Ctx) error {
	var transaction []Entity.Transaction

	err := database.DB.Debug().Find(&transaction).Error
	if err != nil {
		panic(err)
	}

	return ctx.JSON(transaction)
}