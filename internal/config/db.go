package config

import (
	"fmt"
	"log"
	"os"
	"time"

	. "Backend_Go/internal/entities"

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

	db.AutoMigrate(
		&User{},
		&Dealer{},
		&Car{},
		&CarImage{},
		&Lead{},
		&Favorite{},
		&Review{},
		&Report{},
		&RefreshToken{},
	)
	fmt.Println("Config Database Successful..")
	return db, nil

}
