package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title string `json:"title"`
}

type CreateBookPayload struct {
	Title string `json:"title"`
}

type GetBookByIdPayload struct {
	Id uint `json:"id" binding:"required" uri:"id"`
}

type SearchBooksPayload struct {
	TitleQuery string `json:"title_query" form:"title_query"`
}

type UpdateBookPayload struct {
	Id    uint   `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
}

type RemoveBookPayload struct {
	Id uint `json:"id" binding:"required" uri:"id"`
}
