package Controllers

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"net/http"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/Models/Entity"

	"go.mongodb.org/mongo-driver/bson"
)

func GetOneTranscation (c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	transactionId := c.params("transactionId")
	var transaction Entity.Transaction
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(transactionId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&transaction)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.TransactionResponse{
			Status: http.StatusInternalServerError,
			Message: "error",
			Data: &fiber.Map{"data": err.Error()}
		})
	}

	return c.status(http.StatusOK).JSON(response.TransactionResponse{
		Status: http.StatusOK,
		Message : "success",
		Data: &fiber.Map{"data": transaction}
	})
}