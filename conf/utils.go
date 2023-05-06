package conf

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kallepan/ssh-honeypot/logger"
)

func GetValueFromEnv(key string) string {
	value := os.Getenv(key)

	return value
}

func LoadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger.Info("Environment variables loaded")
}
