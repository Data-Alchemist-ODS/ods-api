package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("Database URL is Empty")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoClient = client

	fmt.Println("Successfully Connect To Database...")

	return client
}

// func DisconnectDB() error {
// 	err := mongoClient.Disconnect(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Connection to Database is closed...")

// 	return nil
// }

//client instance

var mongoClient *mongo.Client = ConnectDB()

func GetCollection (client *mongo.Client, name string) *mongo.Collection {
	coll := client.Database("ODS").Collection(name)
	
	return coll
}

