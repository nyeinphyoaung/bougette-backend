package configs

import (
	"bougette-backend/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	ServerIP      string
	ServerPort    string
	DBName        string
	DBHost        string
	DBPort        string
	DBUsername    string
	DBPassword    string
	DB            *gorm.DB
	MAIL_SENDER   string
	MAIL_HOST     string
	MAIL_PORT     string
	MAIL_USERNAME string
	MAIL_PASSWORD string
	VIA_APP_NAME  string
	JWT_SECRET    string
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
	cfg.MAIL_SENDER = getEnvOrDefault("MAIL_SENDER", "")
	cfg.MAIL_HOST = getEnvOrDefault("MAIL_HOST", "")
	cfg.MAIL_PORT = getEnvOrDefault("MAIL_PORT", "2525")
	cfg.MAIL_USERNAME = getEnvOrDefault("MAIL_USERNAME", "")
	cfg.MAIL_PASSWORD = getEnvOrDefault("MAIL_PASSWORD", "")
	cfg.VIA_APP_NAME = getEnvOrDefault("VIA_APP_NAME", "Bougette")
	cfg.JWT_SECRET = getEnvOrDefault("JWT_SECRET", "your_jwt_secret")
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

func (c *Config) InitializedDB() {
	c.DB.AutoMigrate(
		&models.Users{},
		&models.PasswordReset{},
		&models.Categories{},
	)
}

var Envs = loadEnv(".env")
