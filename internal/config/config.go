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
	MongoDBURI        string
	MongoDBName       string
}

func NewConfig() *Config {
	return &Config{
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),
		MongoDBURI:    os.Getenv("MONGODB_URI"),
		MongoDBName:   os.Getenv("MONGODB_NAME"),
	}
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	return NewConfig(), nil
}
