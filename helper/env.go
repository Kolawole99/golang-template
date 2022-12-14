package helper

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnvironmentVariables reads environment variables from a .env file provided or logs an error
func LoadEnvironmentVariables() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
