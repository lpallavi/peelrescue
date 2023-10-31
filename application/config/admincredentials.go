package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AdminName, AdminPw string

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	// Get the api keys for seller and buyer stored as environment variables
	AdminName, _ = os.LookupEnv("ADMIN_NAME")
	AdminPw, _ = os.LookupEnv("ADMIN_PASSWORD")
}
