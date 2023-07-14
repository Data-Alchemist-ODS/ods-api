package controllers

import (
	//default module 
	"context"
	"log"

	//fiber module
	"github.com/gofiber/fiber/v2"

	//mongoDB module
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

	//local module
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
	"github.com/Data-Alchemist-ODS/ods-api/repositories"
)

// TransactionController is a contract what this controller can do
type TransactionController interface {
	//GET HANDLER
	GetAllTransactions(c *fiber.Ctx) error
	GetAllStoredData(c *fiber.Ctx) error

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

//GET REQUEST CONTROLLER
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

//Get All Data From Database
func (controller *transactionController) GetAllStoredData(c *fiber.Ctx) error {
	client := database.GetDB()
	collection := database.GetCollection(client, "Data")

	var dataDocuments []entity.DataDocument

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get data",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &dataDocuments)
	if err != nil {
		log.Fatal(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to decode data",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	var dataResponse []entity.DataResponse

	for _, doc := range dataDocuments {
		for _, data := range doc.Documents {
			dataResponse = append(dataResponse, entity.DataResponse{
				ID:     doc.ID.Hex(),
				Fields: data.Fields,
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Success get all data",
		"status":  fiber.StatusOK,
		"records": dataResponse,
	})
}

//POST REQUEST CONTROLLER
//POST Transaction create by user
func (controller *transactionController) CreateNewTransaction(c *fiber.Ctx) error {
	var request request.TransactionCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	// Save the file data to MongoDB
	err := repositories.SaveToMongoDB(request.FileData)
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