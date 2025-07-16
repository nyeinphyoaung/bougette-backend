package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerIP   string
	ServerPort string
}

func getEnvOrDefault(env string, defaultValue string) string {
	value := os.Getenv(env)
	if value == "" {
		if defaultValue != "" {
			return defaultValue
		} else {
			panic("Environment variable " + env + " is not set")
		}
	}
	return value
}

func LoadEnv(env string) (*Config, error) {
	err := godotenv.Load(env)

	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	cfg.ServerIP = getEnvOrDefault("SERVER_IP", "localhost")
	cfg.ServerPort = getEnvOrDefault("SERVER_PORT", "9984")
	return cfg, nil
}
