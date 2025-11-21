package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv memuat variabel lingkungan dari file .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}