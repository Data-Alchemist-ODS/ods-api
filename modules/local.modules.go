package modules

import (
	//default modules
	"context"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/repositories"
)

// function to store data in data collection mongoDB
func SaveToMongoDB(filename string, c *fiber.Ctx) error {

	coll := database.GetCollection(database.GetDB(), "Data")

	contentType := repositories.GetFileContentType(filename)

	if contentType == "text/csv" {
		file, err := os.Open(filename)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when opening file",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}
		defer file.Close()

		// Read the CSV file
		reader := csv.NewReader(file)
		data, err := reader.ReadAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when reading CSV file",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		headers := data[0]
		documents := make([]map[string]interface{}, 0)

		for i := 1; i < len(data); i++ {
			row := data[i]
			doc := make(map[string]interface{})
			doc["fields"] = make(map[string]interface{}) // Add this line to create the "fields" map

			for j := 0; j < len(headers); j++ {
				fieldValue := row[j]
				fieldType := GetFieldType(fieldValue)

				switch fieldType {
				case "string":
					doc["fields"].(map[string]interface{})[headers[j]] = fieldValue
				case "number":
					numericValue, err := strconv.ParseFloat(fieldValue, 64)
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"message": "error when converting field value to number",
							"status":  fiber.StatusInternalServerError,
							"error":   err.Error(),
						})
					}
					doc["fields"].(map[string]interface{})[headers[j]] = numericValue
				case "boolean":
					booleanValue, err := strconv.ParseBool(fieldValue)
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"message": "error when converting field value to boolean",
							"status":  fiber.StatusInternalServerError,
							"error":   err.Error(),
						})
					}
					doc["fields"].(map[string]interface{})[headers[j]] = booleanValue
				default:
					doc["fields"].(map[string]interface{})[headers[j]] = fieldValue
				}
			}

			documents = append(documents, doc)
		}

		if _, err := coll.InsertOne(context.Background(), bson.M{"documents": documents}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when inserting data to MongoDB",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}
	}

	if contentType == "application/json" {
		fileContent, err := ioutil.ReadFile(filename)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when reading file",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		var jsonData []map[string]interface{}

		err = json.Unmarshal(fileContent, &jsonData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when unmarshaling JSON",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		// Convert JSON data to the desired structure
		documents := make([]map[string]interface{}, len(jsonData))
		for i, doc := range jsonData {
			fields := make(map[string]interface{})
			for key, value := range doc {
				fields[key] = value
			}
			documents[i] = map[string]interface{}{
				"fields": fields,
			}
		}

		if _, err := coll.InsertOne(context.Background(), bson.M{"documents": documents}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when inserting data to MongoDB",
				"status":  fiber.StatusInternalServerError,
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "JSON file inserted successfully",
			"status":  fiber.StatusOK,
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": "error when reading file",
		"status":  fiber.StatusBadRequest,
		"error":   "file type not supported",
	})
}

func GetFieldType(value string) string {
	// Try parsing as boolean
	_, err := strconv.ParseBool(value)
	if err == nil {
		return "boolean"
	}

	// Try parsing as number
	_, err = strconv.ParseFloat(value, 64)
	if err == nil {
		return "number"
	}

	// Default to string if not a boolean or number
	return "string"
}
