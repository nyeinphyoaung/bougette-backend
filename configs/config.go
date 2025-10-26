package configs

import (
	"bougette-backend/models"
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	ServerIP              string
	ServerPort            string
	DBName                string
	DBHost                string
	DBPort                string
	DBUsername            string
	DBPassword            string
	DB                    *gorm.DB
	RedisHost             string
	RedisPort             string
	Redis                 *redis.Client
	MAIL_SENDER           string
	MAIL_HOST             string
	MAIL_PORT             string
	MAIL_USERNAME         string
	MAIL_PASSWORD         string
	VIA_APP_NAME          string
	JWT_SECRET            string
	AWS_REGION            string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_BUCKET_NAME       string
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
	cfg.RedisHost = getEnvOrDefault("REDIS_HOST", "localhost")
	cfg.RedisPort = getEnvOrDefault("REDIS_PORT", "6379")
	cfg.MAIL_SENDER = getEnvOrDefault("MAIL_SENDER", "")
	cfg.MAIL_HOST = getEnvOrDefault("MAIL_HOST", "")
	cfg.MAIL_PORT = getEnvOrDefault("MAIL_PORT", "2525")
	cfg.MAIL_USERNAME = getEnvOrDefault("MAIL_USERNAME", "")
	cfg.MAIL_PASSWORD = getEnvOrDefault("MAIL_PASSWORD", "")
	cfg.VIA_APP_NAME = getEnvOrDefault("VIA_APP_NAME", "Bougette")
	cfg.JWT_SECRET = getEnvOrDefault("JWT_SECRET", "your_jwt_secret")
	cfg.AWS_REGION = getEnvOrDefault("AWS_REGION", "ap-southeast-1")
	cfg.AWS_ACCESS_KEY_ID = getEnvOrDefault("AWS_ACCESS_KEY_ID", "")
	cfg.AWS_SECRET_ACCESS_KEY = getEnvOrDefault("AWS_SECRET_ACCESS_KEY", "")
	cfg.AWS_BUCKET_NAME = getEnvOrDefault("AWS_BUCKET_NAME", "")
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

func (c *Config) ConnectRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisHost + ":" + c.RedisPort,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Successfully connected to Redis")
	c.Redis = rdb
	return nil
}

func (c *Config) InitializedDB() {
	// add drop table
	// c.DB.Migrator().DropTable(
	// 	&models.Users{},
	// 	&models.PasswordReset{},
	// 	&models.Categories{},
	// 	&models.Budgets{},
	// 	&models.Notifications{},
	// 	&models.Wallet{},
	// )

	c.DB.AutoMigrate(
		&models.Users{},
		&models.PasswordReset{},
		&models.Categories{},
		&models.Budgets{},
		&models.Notifications{},
		&models.Wallet{},
	)
}

var Envs = loadEnv(".env")
