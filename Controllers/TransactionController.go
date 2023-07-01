package Controllers

import (
	"context"
	"github.com/gofiber/fiber/v2"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/Models/Entity"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// func GetOneTranscation (c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	transactionId := c.params("transactionId")
// 	var transaction Entity.Transaction
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(transactionId)

// 	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&transaction)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(response.TransactionResponse{
// 			Status: http.StatusInternalServerError,
// 			Message: "error",
// 			Data: &fiber.Map{"data": err.Error()}
// 		})
// 	}

// 	return c.status(http.StatusOK).JSON(response.TransactionResponse{
// 		Status: http.StatusOK,
// 		Message : "success",
// 		Data: &fiber.Map{"data": transaction}
// 	})
// 

func GetAllTransactions(c *fiber.Ctx) error {
	db := database.ConnectDB() // Mengambil koneksi database dari package database
	defer db.Disconnect(context.Background())

	collection := database.GetCollection(db, "Transaction") // Mendapatkan objek koleksi "Transaction"

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
