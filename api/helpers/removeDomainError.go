package helpers

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

// Prevent infinite callback abuse
func RemoveDomainError(url string) bool {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file from helper function")
	}

	genURL := GeneralizeURL(url)
	
	
	if genURL == "http://" + os.Getenv("DOMAIN") {
		return false
	}

	return true
}