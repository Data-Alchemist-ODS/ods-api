package modules

import (
	//default modules
	"encoding/csv"
	"encoding/json"
	"os"
	"context"
	"io/ioutil"	

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
	"github.com/Data-Alchemist-ODS/ods-api/repositories"
)

//function to store data in data collection mongoDB
func SaveToMongoDB(filename string, c *fiber.Ctx) error {

	coll := database.GetCollection(database.GetDB(), "Data")

	file, err := os.Open(filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "error when opening file",
			"status": fiber.StatusInternalServerError,
			"error": err.Error(),
		})
	}
	defer file.Close()

	contentType := repositories.GetFileContentType(filename)
	if contentType == "text/csv" {
		// Read the CSV file
		reader := csv.NewReader(file)
		data, err := reader.ReadAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when reading csv file",
				"status": fiber.StatusInternalServerError,
				"error": err.Error(),
			})
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when inserting data to mongoDB",
				"status": fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}
	}
	
	//STILL GOT AN ERROR HERE
	if contentType == "application/json" {
		filePath := file.Name()

		data, err := ioutil.ReadFile(filePath)
		if err != nil{
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when opening file",
				"status": fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		var jsonData []map[string]interface{}
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when unmarshalling json",
				"status": fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}

		columns := make([]string, 0)
		for key := range jsonData[0] {
			columns = append(columns, key)
		}

		content := make([][]string, 0)
		content = append(content, columns)

		for _, item := range jsonData {
			row := make([]string, 0)
			for _, col := range columns {
				value, ok := item[col].(string)
				if !ok {
					continue
				}
				row = append(row, value)
			}

			content = append(content, row)
		}

		documents := make([]request.Data, 0)
		headers := content[0]
		for i := 1; i < len(content); i++ {
			row := content[i]
			doc := request.Data{
				Fields: make(map[string]string),
			}

			for j := 0; j < len(headers); j++ {
				doc.Fields[headers[j]] = row[j]
			}

			documents = append(documents, doc)
		}

		if _, err := coll.InsertOne(context.Background(), bson.M{"documents": documents}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "error when inserting data to mongoDB",
				"status": fiber.StatusInternalServerError,
				"error": err.Error(),
			})
		}
	}
	return nil
}