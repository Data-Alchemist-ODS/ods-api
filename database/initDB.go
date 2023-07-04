package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Data-Alchemist-ODS/ods-api/config"
)

var mongoClient *mongo.Client

func ConnectDB() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.LoadENV()))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Connect To Database...")

	mongoClient = client 

	return client
}

func DisconnectDB() error {
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to Database is closed...")

	return nil
}

func GetCollection(client *mongo.Client, name string) *mongo.Collection { // Mengubah nama fungsi menjadi GetCollection
	coll := client.Database("ODS").Collection(name)

	return coll
}

func GetDB() *mongo.Client {
	return mongoClient
}

func LoadMongoURI() string {
	URI := godotenv.Load()
	if URI != nil {
		log.Fatal("Error loading MongoURI file")
	}

	return URI
}
