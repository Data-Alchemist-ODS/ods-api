package Controllers

import (
	"time"
	"context"
	"github.com/gofiber/fiber/v2"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/Models/Entity"
	"github.com/Data-Alchemist-ODS/ods-api/Models/Request"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllTransactions(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB() // Mengambil koneksi database dari package database

	collection := database.GetCollection(client, "Transaction") // Mendapatkan objek koleksi "Transaction"

	var transactions []Entity.Transaction
	cursor, err := collection.Find(context.Background(), options.Find())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get transactions",
			"error":   err.Error(),
		})
	}

	if err := cursor.All(context.Background(), &transactions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to decode transactions",
			"error":   err.Error(),
		})
	}

	return c.JSON(transactions)
}

func CreateTransaction(c *fiber.Ctx) error {
    db := database.ConnectDB()
    defer db.Disconnect(context.Background())

    var request Request.TransactionCreateRequest
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "invalid request body",
            "error":   err.Error(),
        })
    }

	var records [][]string
	var err error
	if request.FileContentType == "text/csv" {
		records, err = repositories.ReadCSV(request.FileData)
	} else if request.FileContentType == "application/json" {
		jsonData, err := repositories.ReadJSON(request.FileData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to read json file",
				"error":   err.Error(),
			})
		}

		records = repositories.ConvertJSONToCSV(jsonData)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid file format",
			"error":   err.Error(),
		})
	}

    transaction := Entity.Transaction{
        PartitionType: request.PartitionType,
        ShardingKey:   request.ShardingKey,
        Database:      request.Database,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
		Data: record,
    }

    collection := database.GetCollection(db, "Transaction")
    _, err := collection.InsertOne(context.Background(), transaction)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "failed to create transaction",
            "error":   err.Error(),
        })
    }

    return c.JSON(fiber.Map{
		"transaction": transaction,
		"message": "transaction created successfully"})
}
