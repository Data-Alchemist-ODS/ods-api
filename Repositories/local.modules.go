package repositories

import (
	//default modules
	"encoding/csv"
	"os"
	"context"
	"log"

	//mongoDB modules
	"go.mongodb.org/mongo-driver/bson"

	//local modules
	"github.com/Data-Alchemist-ODS/ods-api/database"
	"github.com/Data-Alchemist-ODS/ods-api/models/request"
)

//function to store data in data collection mongoDB
func SaveToMongoDB(FileData string) error {

	coll := database.GetCollection(database.GetDB(), "Data")

	file, err := os.Open(FileData)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
		return err
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
		log.Fatal(err)
		return err
	}

	return nil
}