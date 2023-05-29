package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalsqwe/go-api/config/db"
	"github.com/mikalsqwe/go-api/config/env"
	"github.com/mikalsqwe/go-api/models"
	"github.com/mikalsqwe/go-api/utils"
)

func main() {
	r := SetupServer(env.DATABASE_URL)

	r.Run(env.PORT)
}

func CreateBook(c *gin.Context) {
	var payload models.CreateBookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error,
		})
		return
	}

	book := models.Book{
		Title: payload.Title,
	}

	if err := db.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
	return
}

func GetBookById(c *gin.Context) {
	var payload models.GetBookByIdPayload

	if err := c.BindUri(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error,
		})
		return
	}

	var book models.Book

	if err := db.DB.Where("id = ?", payload.Id).First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
	return
}

func GetBooksPagination(c *gin.Context) {
	var payload utils.PaginationPayload

	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   err.Error,
			"message": "payload error",
		})
		return
	}

	var books []models.Book

	pagination := utils.Pagination[[]models.Book]{
		Page:   payload.Page,
		Limit:  payload.Limit,
		SortBy: payload.SortBy,
	}

	if err := db.DB.Scopes(utils.Paginate(books, &pagination, db.DB)).Find(&books).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error,
			"message": "scopes error",
		})
		return
	}

	pagination.Rows = books

	c.JSON(http.StatusOK, gin.H{
		"data": pagination,
	})
	return
}

func SearchBooks(c *gin.Context) {
	var payload models.SearchBooksPayload

	if err := c.Bind(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error,
		})
		return
	}

	var books []models.Book

	interpolatedTitleQuery := "%" + payload.TitleQuery + "%"

	if err := db.DB.Where("title LIKE ?", interpolatedTitleQuery).Find(&books).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": books,
	})
	return
}

func UpdateBook(c *gin.Context) {
	var payload models.UpdateBookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error,
			"message": "payload error",
		})
		return
	}

	var book models.Book

	if err := db.DB.Where("id = ?", payload.Id).Find(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error,
			"message": "db find error",
		})
		return
	}

	book.Title = payload.Title

	if err := db.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error,
			"message": "db save error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
	return
}

func RemoveBook(c *gin.Context) {
	var payload models.RemoveBookPayload

	if err := c.BindUri(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error,
			"message": "payload error",
		})
		return
	}

	var book models.Book

	if err := db.DB.Where("id = ?", payload.Id).Delete(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   err.Error,
			"message": "detele error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": book,
	})
	return
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/books", GetBooksPagination)
	r.GET("/books/:id", GetBookById)
	r.GET("/books/search", SearchBooks)
	r.POST("/books/create", CreateBook)
	r.PATCH("/books/update", UpdateBook)
	r.DELETE("/books/remove/:id", RemoveBook)

	return r
}

func SetupServer(dbName string) *gin.Engine {
	db.ConnectDatabase(dbName)

	r := setupRouter()

	return r
}
