package db

import (
	"log"

	"github.com/mikalsqwe/go-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dbName string) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed connecting to the database: %v", dbName)
	}

	db.AutoMigrate(&models.Book{})

	DB = db
}
