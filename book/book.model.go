package book

import (
	"github.com/mikalsqwe/go-api/author"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string          `json:"title" binding:"required"`
	Description string          `json:"description" binding:"required"`
	AuthorsIds  []uint          `json:"authorsIds"`
	Authors     []author.Author `json:"authors"`
}

type CreateBookPayload struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	AuthorsIds  []string `json:"authorsIds"`
}

type GetBookByIdPayload struct {
	Id string `json:"id"`
}
