package db

import (
	"log"

	"github.com/mikalsqwe/go-api/author"
	"github.com/mikalsqwe/go-api/book"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dbName string) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed connecting to the database: %v", dbName)
	}

	db.AutoMigrate(&book.Book{}, &author.Author{})
}
