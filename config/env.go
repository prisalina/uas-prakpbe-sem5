package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found in working dir; using system environment")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
