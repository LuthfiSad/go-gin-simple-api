package config

import (
	"fmt"
	"go-gin-simple-api/model"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	ServerPort       string
	JWTSecret        string
	JWTExpiryHours   string
	CloudinaryName   string
	CloudinaryKey    string
	CloudinarySecret string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	config := &Config{
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		ServerPort:       os.Getenv("SERVER_PORT"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTExpiryHours:   os.Getenv("JWT_EXPIRATION"),
		CloudinaryName:   os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinarySecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}

	return config, nil
}

func SetupDatabase(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Lakukan AutoMigrate untuk semua model yang kamu pakai
	err = db.AutoMigrate(
		&model.Book{},
		&model.Media{},
		&model.User{},
		&model.BookStock{},
		&model.BookTransaction{},
		&model.Customer{},
		&model.Charge{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
