package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func LoadDB() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	sqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", host, user, password, dbName)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println("Connected to database")

	DB = db
}
