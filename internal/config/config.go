// internal/config/config.go
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret         string
	JWTExpiration     string
	GoogleCredentials string
	GPTAPI1Key        string
	GPTAPI2Key        string
	MongoDBURI        string
	MongoDBName       string
}

func NewConfig() *Config {
	return &Config{
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTExpiration:     os.Getenv("JWT_EXPIRATION"),
		GoogleCredentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		GPTAPI1Key:        os.Getenv("GPTT_APII_KEYY_1"),
		GPTAPI2Key:        os.Getenv("GPTT_APII_KEYY_2"),
		MongoDBURI:        os.Getenv("MONGODB_URI"),
		MongoDBName:       os.Getenv("MONGODB_NAME"),
	}
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	return NewConfig(), nil
}
