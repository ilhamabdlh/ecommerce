package config

import (
	"log"
	"os"

	"github.com/ilhamabdlh/ecommerce/pkg/database"
)

type Config struct {
	MongoDB   *database.MongoDB
	JWTSecret string
}

func NewConfig() *Config {
	mongoURI := getEnv("MONGODB_URI", "mongodb://localhost:27017")
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")

	mongodb, err := database.NewMongoDB(mongoURI)
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		MongoDB:   mongodb,
		JWTSecret: jwtSecret,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
