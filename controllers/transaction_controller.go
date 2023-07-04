package controllers

import (
	"context"
	"encoding/csv"
	"io/ioutil"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/Data-Alchemist-ODS/ods-api/config"
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
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

func saveFileData(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

type Person struct {
	gorm.Model
	Fields map[string]string `gorm:"-"`
}

func SaveToMongoDB(PartitionType, ShardingKey, Database, FileData string) error {
	db := database.ConnectDB()

	dsn := config.LoadENV()

	db, err := gorm.Open(mongo.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	defer db.Disconnect(context.Background())

	db.AutoMigrate(&Person{})

	file, err := os.Open(FileData)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, row := range data {
		history := Person{
			Fields: make(map[string]string),
		}
		for i := 0; i < len(row); i++ {
			fieldName := "Field" + string(i+1)
			history.Fields[fieldName] = row[i]
		}

		if err := db.Create(&history).Error; err != nil {
			return err
		}
	}

	return nil
}

func (controller *transactionController) CreateTransaction(c *fiber.Ctx) error {
	var request request.TransactionCreateRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	err := SaveToMongoDB(request.PartitionType, request.ShardingKey, request.Database, request.FileData)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "data has been saved to mongo",
	})
}
