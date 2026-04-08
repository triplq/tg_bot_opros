package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Sprintf("Error loading .env file: %v", err)
	}
	value, ok := os.LookupEnv(key)
	if !ok {
		return fmt.Sprintf("environment variable %q is not set", key)
	}

	return value
}
