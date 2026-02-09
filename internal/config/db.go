package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"Backend_Go/internal/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (db *gorm.DB, err error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	// Auto Migrate
	err = db.AutoMigrate(
		&entities.User{},
		&entities.Dealer{},
		&entities.Car{},
		&entities.CarImage{},
		&entities.Lead{},
		&entities.Favorite{},
		&entities.Review{},
		&entities.Report{},
		&entities.RefreshToken{},
		&entities.Conversation{},
		&entities.Message{},
	)
	if err != nil {
		return nil, err
	}

	// MIGRATION: Fix existing cars with empty status -> 'approved'
	// This ensures existing cars don't disappear from public listing.
	// Only affects rows where status is NULL or empty string.
	if err := db.Model(&entities.Car{}).Where("status IS NULL OR status = ''").Update("status", "approved").Error; err != nil {
		log.Printf("Migration warning: failed to update car statuses: %v", err)
	} else {
		log.Println("Migration: updated empty car statuses to 'approved'")
	}

	fmt.Println("Config Database Successful..")
	return db, nil

}
