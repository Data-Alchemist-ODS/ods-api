package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func connectDB () error {
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

func disconnectDB () error {
	err := mongoClient.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}	

	fmt.Println("Connection to Database is closed...")

	return nil
}

func main() {
	// Read env
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	connectDB()

	//sleep time 5 second
	time.Sleep(5 * time.Second)

	// Disconnect to database
	disconnectDB()
}
