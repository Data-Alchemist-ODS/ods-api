package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
	"github.com/Data-Alchemist-ODS/ods-api/repositories"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionController is a contract what this controller can do
type TransactionController interface {
	GetAllTransactions(c *fiber.Ctx) error
	CreateTransaction(c *fiber.Ctx) error
}

// transactionController is a struct that represent the TransactionController contract
type transactionController struct{}

// NewTransactionController is the constructor
func NewTransactionController() TransactionController {
	return &transactionController{}
}

/*
 *  Implement functions goes down here
 */

func (controller *transactionController) GetAllTransactions(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB() // Mengambil koneksi database dari package database

	collection := database.GetCollection(client, "Transaction") // Mendapatkan objek koleksi "Transaction"

	var transactions []entity.Transaction
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

func (controller *transactionController) CreateTransaction(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	var request request.TransactionCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	fileContentType := repositories.GetFileContentType(request.FileName)

	records, err := repositories.ReadData(request.FileData, fileContentType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to read data",
			"error":   err.Error(),
		})
	}

	transaction := entity.Transaction{
		PartitionType: request.PartitionType,
		ShardingKey:   request.ShardingKey,
		Database:      request.Database,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Data:          records,
	}

	collection := database.GetCollection(db, "Transaction")
	_, err = collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to create transaction",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"transaction": transaction,
		"message":     "transaction created successfully",
	})
}
