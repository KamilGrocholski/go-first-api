package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalsqwe/go-api/book"
	"github.com/mikalsqwe/go-api/db"
	"github.com/mikalsqwe/go-api/utils"
)

func main() {
	Start("main.db")
}

func Start(dbName string) {
	db.ConnectDatabase(dbName)

	r := gin.Default()

	r.POST("/create/book", func(c *gin.Context) {
		var payload book.CreateBookPayload

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": ";/",
			})
		}

		authorsIds := make([]uint, 0)

		for index, dirtyAuthorId := range payload.AuthorsIds {
			authorId, err := utils.StringToUint(dirtyAuthorId)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": ";/",
				})
				return
			}

			authorsIds[index] = authorId
		}

		book := &book.Book{
			Title:       payload.Title,
			Description: payload.Description,
			AuthorsIds:  authorsIds,
		}

		err := db.DB.Create(book).Error

		if err == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "db error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": book,
		})
	})

	r.GET("/book:id", func(c *gin.Context) {
		dirtyId := c.Params.ByName("id")
		bookId, idParsingErr := utils.StringToUint(dirtyId)

		if idParsingErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "parsing error",
			})
			return
		}

		var book book.Book

		err := db.DB.First(book, bookId).Error

		if err == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": book,
		})
	})

	r.Run()
}
