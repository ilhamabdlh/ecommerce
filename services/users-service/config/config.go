package config

type Config struct {
	MongoDB struct {
		URI      string
		Database string
	}
	Server struct {
		Port string
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// Default values
	cfg.MongoDB.URI = "mongodb://localhost:27017"
	cfg.MongoDB.Database = "user_service"
	cfg.Server.Port = ":8081"

	return cfg, nil
}
