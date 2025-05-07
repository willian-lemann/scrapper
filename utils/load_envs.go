package utils

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %s", err)
	}
}
