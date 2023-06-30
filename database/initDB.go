package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

var DB *gorm.DB
var mongoClient *mongo.Client

func GetCollection(name string) *mongo.Collection {
	dbName := os.Getenv("DATABASE")
	if dbName == "" {
		log.Fatal("Database Name Is Empty")
	}

	col := mongoClient.Database(dbName).Collection(name)

	fmt.Println("Successfully Get Collection:", col.Name())
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
