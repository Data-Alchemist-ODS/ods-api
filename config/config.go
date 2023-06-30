package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//function config
func LoadENV () string {
	err := godotenv.Load()
	err != nil {
		log.Fatal("Error loading .env file")
	}

	//return os.getENV
	return os.Getenv("DATABASE_URL")
}