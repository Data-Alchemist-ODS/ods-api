package controllers

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
	"github.com/Data-Alchemist-ODS/ods-api/repositories"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllTransactions(c *fiber.Ctx) error {
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

func CreateTransaction(c *fiber.Ctx) error {
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

func ConnectTiDB(c *fiber.Ctx) error {
	// get from request body
	serverName := c.FormValue("server_name")
	user := c.FormValue("user")
	password := c.FormValue("password")
	database := c.FormValue("database")

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		// ServerName: "gateway01.eu-central-1.prod.aws.tidbcloud.com",
		ServerName: serverName,
	})

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:4000)/%s?tls=tidb", user, password, serverName, database)
	// db, err := sql.Open("mysql", "4MXeBRmXXzc7uqt.root:<your_password>@tcp(gateway01.eu-central-1.prod.aws.tidbcloud.com:4000)/test?tls=tidb")
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to connect to TiDB",
			"error":   err.Error(),
		})
	}

	// TODO save to cache
	// ...

	return c.JSON(fiber.Map{
		"message": "connected to TiDB",
		"db":      db,
	})
}
