package database

import (
	"log"

	"github.com/crlnravel/go-fiber-template/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgres() {
	dsn := config.GetEnv("DATABASE_URL", "")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db

	log.Print("Successfully connected to database")
}
