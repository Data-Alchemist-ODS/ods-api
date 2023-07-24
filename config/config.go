package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//Database function config
func LoadENV() string {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("error loading .env file")
	}

	//return os.getENV
	return os.Getenv("DATABASE_URL")
}

//API function config
func LoadAPIKey() string {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("error loading .env file")
	}

	//return os.getENVq
	return os.Getenv("API_KEY")
}