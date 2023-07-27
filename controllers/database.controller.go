package controllers

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Data-Alchemist-ODS/ods-api/models/request"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

// DatabaseController is a contract what this controller can do
type DatabaseController interface {
	ConnectToTiDB(c *fiber.Ctx) error
}

// databaseController is a struct that represent the DatabaseController contract
type databaseController struct{}

// NewDatabaseController is the constructor
func NewDatabaseController() DatabaseController {
	return &databaseController{}
}

/*
 *  Implement functions goes down here
 */

func (controller *databaseController) ConnectToTiDB(c *fiber.Ctx) error {
	var request request.TiDBConnectionRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "failed to parse json",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		// ServerName: "gateway01.eu-central-1.prod.aws.tidbcloud.com",
		ServerName: request.ServerName,
	})

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:4000)/%s?tls=tidb", request.User, request.Password, request.ServerName, request.Database)
	// db, err := sql.Open("mysql", "4MXeBRmXXzc7uqt.root:<your_password>@tcp(gateway01.eu-central-1.prod.aws.tidbcloud.com:4000)/test?tls=tidb")
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		log.Fatal("failed to connect database", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to connect to TiDB",
			"error":   err.Error(),
		})
	}

	defer db.Close()

	ca := cache.New(5*time.Minute, 10*time.Minute)

	// Set the value of the key "foo" to "bar", with the default expiration time
	ca.Set("tidb_connection", dataSourceName, cache.NoExpiration)

	connection, found := ca.Get("tidb_connection")

	if !found {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to connect to TiDB",
			"error":   err.Error(),
		})
	}

	var dbName string
	if err = db.QueryRow("SELECT * FROM `fortune500_2018_2022`").Scan(&dbName); err != nil {
		log.Fatal("failed to execute query", err)
	}
	fmt.Println(dbName)

	return c.JSON(fiber.Map{
		"message":    "connected to TiDB",
		"connection": connection,
		"db":         db,
	})
}
