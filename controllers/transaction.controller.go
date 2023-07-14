package controllers

import (
	"context"
	"encoding/csv"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// TransactionController is a contract what this controller can do
type TransactionController interface {
	//GET HANDLER
	GetAllTransactions(c *fiber.Ctx) error
	// GetAllStoredData(c *fiber.Ctx) error

	//POST HANDLER
	CreateNewTransaction(c *fiber.Ctx) error

	//UPDATE HANDLER
	//DELETE HANDLER
}

// transactionController is a struct that represent the TransactionController contract
type transactionController struct{}

// NewTransactionController is the constructor
func NewTransactionController() TransactionController {
	return &transactionController{}
}

//Get All Transaction Done By User
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
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	if err := cursor.All(context.Background(), &transactions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to decode transactions",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success get all transactions",
		"status":  fiber.StatusOK,
		"records": transactions,
	})
}

//Get One Data From Database
// func (controller *transactionController) GetAllStoredData(c *fiber.Ctx) {
// 	db := database.ConnectDB()
// 	defer db.Disconnect(context.Background())

// 	client := database.GetDB()

// 	collection := database.GetCollection(client, "Data")
// }

func SaveToMongoDB(FileData string) error {

	coll := database.GetCollection(database.GetDB(), "Data")

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

	documents := make([]request.Data, 0) // Changed the type to []Data

	headers := data[0]
	for i := 1; i < len(data); i++ {
		row := data[i]
		doc := request.Data{
			Fields: make(map[string]string),
		}

		for j := 0; j < len(headers); j++ {
			doc.Fields[headers[j]] = row[j]
		}

		documents = append(documents, doc)
	}

	if _, err := coll.InsertOne(context.Background(), bson.M{"documents": documents}); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

//POST Transaction create by user
func (controller *transactionController) CreateNewTransaction(c *fiber.Ctx) error {
	var request request.TransactionCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	// Save the file data to MongoDB
	err := SaveToMongoDB(request.FileData)
	if err != nil {
		return err
	}

	collection := database.GetCollection(database.GetDB(), "Transaction")

	transaction := entity.Transaction{
		PartitionType: request.PartitionType,
		ShardingKey:   request.ShardingKey,
		Database:      request.Database,
		Data:          request.FileData,
	}
	transaction.ID = primitive.NewObjectID()

	_, err = collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create transaction",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Success create transaction",
		"status":  fiber.StatusOK,
		"record":  transaction,
	})
}