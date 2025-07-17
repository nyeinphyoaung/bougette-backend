package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	ServerIP   string
	ServerPort string
	DBName     string
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	DB         *gorm.DB
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

func loadEnv(env string) *Config {
	godotenv.Load(env)

	cfg := &Config{}
	cfg.ServerIP = getEnvOrDefault("SERVER_IP", "localhost")
	cfg.ServerPort = getEnvOrDefault("SERVER_PORT", "9984")
	cfg.DBName = getEnvOrDefault("DB_NAME", "bougette")
	cfg.DBHost = getEnvOrDefault("DB_HOST", "localhost")
	cfg.DBPort = getEnvOrDefault("DB_PORT", "3306")
	cfg.DBUsername = getEnvOrDefault("DB_USERNAME", "nyeinphyoaung")
	cfg.DBPassword = getEnvOrDefault("DB_PASSWORD", "password")
	return cfg
}

func (c *Config) ConnectDB() error {
	dsn := c.DBUsername + ":" + c.DBPassword + "@tcp(" + c.DBHost + ")/" + c.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	log.Println("Successfully connected to database")
	c.DB = db
	return nil
}

var Envs = loadEnv(".env")
