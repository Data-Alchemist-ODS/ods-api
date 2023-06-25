package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("There's Problem Loading .env File")
	}

	DB_URL := os.Getenv("DATABASE_URL")
	if DB_URL == "" {
		log.Fatal("Database URL is Empty")
	}

	fmt.Println("My Database URL is:", DB_URL)

	clientoption := options.Client().ApplyURI(DB_URL)

	client, err := mongo.Connect(context.Background(), clientoption)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully Connect To Database...")

	err = client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to Database is closed...")

	return client
}
