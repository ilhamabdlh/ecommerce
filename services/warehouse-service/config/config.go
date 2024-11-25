package config

import (
	"os"
)

type Config struct {
	MongoURI               string
	DatabaseName           string
	ServerPort             string
	JWTSecret              string
	ConsulAddr             string
	ServiceName            string
	ServiceID              string
	EnableServiceDiscovery bool
}

func LoadConfig() *Config {
	config := &Config{
		MongoURI:               getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:           getEnv("DB_NAME", "warehouse_db"),
		ServerPort:             getEnv("SERVER_PORT", "8084"),
		JWTSecret:              getEnv("JWT_SECRET", "your-secret-key"),
		ConsulAddr:             getEnv("CONSUL_ADDR", "localhost:8500"),
		ServiceName:            getEnv("SERVICE_NAME", "warehouse-service"),
		ServiceID:              getEnv("SERVICE_ID", "warehouse-1"),
		EnableServiceDiscovery: getEnvBool("ENABLE_SERVICE_DISCOVERY", false),
	}
	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}
