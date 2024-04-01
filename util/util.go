package util

import (
	"os"

	"github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		print("Error loading .env file")
	}

	return os.Getenv(key)
}
