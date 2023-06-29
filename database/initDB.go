package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func GetCollection(name string) *mongo.Collection {
	dbName, err := os.Getenv("DATABASE")
	if err != nil {
		log.Fatal(err)
	}
	
	col := mongoClient.Database(dbName).Collection(name)
	fmt.Println(dbName)
	fmt.Println("Successfully Get Collection...")
	return col
}

func ConnectDB() error {
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

	return nil
}

func DisconnectDB() error {
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to Database is closed...")

	return nil
}
