package modules

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
)

// function to store data in data collection mongoDB
func SaveToMongoDB(filename string, c *fiber.Ctx) error {

	coll := database.GetCollection(database.GetDB(), "Data")

    file, err := c.FormFile(filename)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "error when opening file using form file",
            "status":  fiber.StatusInternalServerError,
            "error":   err.Error(),
        })
    }

    content, err := file.Open()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "error when opening file",
            "status":  fiber.StatusInternalServerError,
            "error":   err.Error(),
        })
    }
    defer content.Close() // Close the opened file instead of the *multipart.FileHeader

	// Determine the file type based on its content type
	contentType := file.Header.Get("Content-Type")

	if contentType == "text/csv" {
		// Read the CSV file
		reader := csv.NewReader(content)
		var headers []string
		var data []request.Data

		for i := 0; ; i++ {
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "error when reading CSV file",
					"status":  fiber.StatusInternalServerError,
					"error":   err.Error(),
				})
			}

			if i == 0 {
				headers = row
				continue
			}

			child := make(request.Data)
			for j, value := range row {
				child[headers[j]] = value
			}

			data = append(data, child)
		}

		if _, err := coll.InsertOne(context.Background(), bson.M{"documents": filename, "data": data}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when inserting data to MongoDB",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}
	} else if contentType == "application/json" {
		fileContent, err := ioutil.ReadAll(content)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when reading file",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		var data []request.Data

		err = json.Unmarshal(fileContent, &data)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when unmarshaling JSON",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		if _, err := coll.InsertOne(context.Background(), bson.M{"documents": filename, "data": data}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when inserting data to MongoDB",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "file type not supported",
			"status":  fiber.StatusBadRequest,
			"error":   "file must be .csv or .json",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "File inserted successfully",
		"status":  fiber.StatusOK,
	})
}
