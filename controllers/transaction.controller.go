package controllers

import (
	//default modules
	"context"
	"fmt"
	"strings"

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/helpers"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
	"github.com/Data-Alchemist-ODS/ods-api/modules"
	"github.com/Data-Alchemist-ODS/ods-api/repositories"

	//third party modules
	"github.com/go-playground/validator/v10"
)

// validator is a variable that represent the validator module
var validate = validator.New()

// TransactionController is a contract what this controller can do
type TransactionController interface {
	//GET HANDLER
	GetAllTransactions(c *fiber.Ctx) error
	GetOneTransaction(c *fiber.Ctx) error

	GetAllStoredDatas(c *fiber.Ctx) error
	GetOneStoredData(c *fiber.Ctx) error

	//POST HANDLER
	CreateNewTransaction(c *fiber.Ctx) error

	//DELETE HANDLER
	DeleteTransaction(c *fiber.Ctx) error
}

// transactionController is a struct that represent the TransactionController contract
type transactionController struct{}

// NewTransactionController is the constructor
func NewTransactionController() TransactionController {
	return &transactionController{}
}

// GET REQUEST CONTROLLER
// Get All Transaction Done By User
func (controller *transactionController) GetAllTransactions(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB() // Mengambil koneksi database dari package database
	collection := database.GetCollection(client, "Transaction") // Mendapatkan objek koleksi "Transaction"

	var transactions []entity.Transaction

	cursor, err := collection.Find(context.Background(), options.Find())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get transactions",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	if err := cursor.All(context.Background(), &transactions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to decode transactions",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success get all transactions",
		"status":  fiber.StatusOK,
		"records": transactions,
	})
}

// Get One Transaction By Id Params
func (controller *transactionController) GetOneTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id format",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "Transaction")

	var transaction entity.Transaction

	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "transaction not found in document",
				"status":  fiber.StatusNotFound,
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get transaction",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "success get transaction",
		"status":  fiber.StatusOK,
		"record":  transaction,
	})
}

// Get All Data From Database
func (controller *transactionController) GetAllStoredDatas(c *fiber.Ctx) error {
	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "Data")

	var results []entity.Document

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to get data",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	err = cursor.All(context.Background(), &results)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to decode data",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "success get all data",
		"status":  fiber.StatusOK,
		"records": results,
	})
}

// Get One Data From Database By Id Params
func (controller *transactionController) GetOneStoredData(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id format",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "Data")

	var results entity.Document
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&results)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Data not found",
				"status":  fiber.StatusNotFound,
			})
		}
	}
	return c.JSON(fiber.Map{
		"message": "Success get data by ID",
		"status":  fiber.StatusOK,
		"records": results,
	})
}

// POST REQUEST CONTROLLER
// POST Transaction create by user
func (controller *transactionController) CreateNewTransaction(c *fiber.Ctx) error {
	var request request.TransactionCreateRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse json",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	if err := validate.Struct(request); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "invalid " + err.Field() + " format",
				"status":  fiber.StatusBadRequest,
				"error":   err.Error(),
			})
		}
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	//Perform Sharding Using Local Modules In Repositories
	if request.PartitionType == "Horizontal" {

		method, err := repositories.HorizontalSharding("data", request.ShardingKey, request.Database, c)

		fmt.Println(len(method))

		helpers.SaveToTiDB(method[0], "gateway01.eu-central-1.prod.aws.tidbcloud.com", "4MXeBRmXXzc7uqt.root", "NLRxAVAVtAKY5SXu", "fortune500")

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to perform sharding",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		// Save the file data to MongoDB
		err = modules.SaveToMongoDB("data", c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to save data to MongoDB",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		file, err := c.FormFile("data")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when opening file",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		fileData := entity.FileData{
			FileName: file.Filename,
		}

		client := database.GetDB()
		collection := database.GetCollection(client, "Transaction")

		transaction := entity.Transaction{
			PartitionType: request.PartitionType,
			ShardingKey:   request.ShardingKey,
			Database:      strings.Join(request.Database, ","),
			Data:          fileData,
		}
		transaction.ID = primitive.NewObjectID()

		_, err = collection.InsertOne(context.Background(), transaction)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "failed to create transaction",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "success create transaction",
			"status":  fiber.StatusOK,
			"record":  transaction,
		})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": "failed to create transaction",
		"status":  fiber.StatusInternalServerError,
		"error":   "invalid partition type",
	})
}

// DELETE REQUEST CONTROLLER
// Delete Transaction By Id Params
func (controller *transactionController) DeleteTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid id format",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	db := database.ConnectDB()
	defer db.Disconnect(context.Background())

	client := database.GetDB()
	collection := database.GetCollection(client, "Transaction")

	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to delete transaction",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "transaction not found",
			"status":  fiber.StatusNotFound,
		})
	}

	return c.JSON(fiber.Map{
		"message": "success delete transaction",
		"status":  fiber.StatusOK,
	})
}
