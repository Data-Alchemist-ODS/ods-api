package controllers

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	if err = rdb.Set(context.Background(), "tidb_connection", dataSourceName, 0).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "failed to save to cache",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "connected to TiDB",
		"db":      db,
	})
}
