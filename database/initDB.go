package database

import (
	"context"
	"fmt"
	"log"
	"os"

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

	mongoClient = client // Assign the client to the global variable

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

func GetColletion(client *mongo.Client, name string) *mongo.Collection {
	coll := client.Database("ODS").Collection(name)

	return coll
}
