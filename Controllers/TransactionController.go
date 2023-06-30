package Controllers

import (
	"encoding/csv"
	"encoding/json"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cespare/xxhash/v2"
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