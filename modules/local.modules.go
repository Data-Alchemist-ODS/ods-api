package modules

import (
	//default modules
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"
	"errors"
	"fmt"
	"strings"

	//fiber modules
	"github.com/gofiber/fiber/v2"

	//mongo modules
	"go.mongodb.org/mongo-driver/bson"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
	"github.com/Data-Alchemist-ODS/ods-api/models/entity"
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

	fileData := entity.FileData{
		FileName: file.Filename,
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

		if _, err := coll.InsertOne(context.Background(), bson.M{"documents": fileData, "data": data}); err != nil {
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

		if _, err := coll.InsertOne(context.Background(), bson.M{"documents": fileData, "data": data}); err != nil {
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

func ConvertToUSADate(indoDate string) (string, error) {
    // Mapping of Indonesian months to English months
    monthMap := map[string]string{
        "Januari":   "January",
        "Februari":  "February",
        "Maret":     "March",
        "April":     "April",
        "Mei":       "May",
        "Juni":      "June",
        "Juli":      "July",
        "Agustus":   "August",
        "September": "September",
        "Oktober":   "October",
        "November":  "November",
        "Desember":  "December",
    }

    // Split the input date and get the month and year
    parts := strings.Fields(indoDate)
    if len(parts) != 3 {
        return "", errors.New("invalid date format")
    }
    day := parts[0]
    month := parts[1]
    year := parts[2]

    // Convert the month to English using the mapping
    englishMonth, found := monthMap[month]
    if !found {
        return "", errors.New("invalid month")
    }

    // Combine the parts and parse the date
    englishDate := fmt.Sprintf("%s %s %s", day, englishMonth, year)
    t, err := time.Parse("2 January 2006", englishDate)
    if err != nil {
        return "", err
    }

    // Format the time.Time in USA date format (month/day/year)
    usaDate := t.Format("1/2/2006")
    return usaDate, nil
}
