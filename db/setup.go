package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"), os.Getenv("DB_TIMEZONE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		log.Fatal("Failed to connecting database: ", err)
	}

	pgDb, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance: ", err)
	}

	pgDb.SetMaxOpenConns(5)
	pgDb.SetMaxIdleConns(2)
	pgDb.SetConnMaxLifetime(time.Hour)

	DB = db.Debug()
	// DB = db

}
