package main

import (
	"net/http"
	"strconv"

	"github.com/mikalsqwe/go-api/lru_cache"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/cos", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	return r
}

func setupDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func Start(db *gorm.DB) {
	moviesCache := lru_cache.NewLruCache[uint, Movie](1000)

	// Migrate the schema
	db.AutoMigrate(&Movie{})

	r := gin.Default()

	r.GET("/movies", func(c *gin.Context) {
		movies := []Movie{}
		db.Find(&movies)

		c.JSON(http.StatusOK, gin.H{
			"movies": movies,
		})
	})

	r.GET("/movie", func(c *gin.Context) {
		queryId := c.Query("id")
		id, success := StringToUint(queryId)

		if success == false {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid id",
			})
		}

		movie, is := GetMovieById(id, db, moviesCache)

		if is == true {
			c.JSON(http.StatusOK, gin.H{
				"movie": movie,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "not found",
			})
		}
	})

	r.POST("/movie/create", func(c *gin.Context) {
		title := c.Query("title")

		movie := CreateMovie(title)

		db.Create(movie)
	})

	r.Run()
}

func main() {
	db := setupDatabase()

	Start(db)
}

type Movie struct {
	gorm.Model
	Title string
}

func CreateMovie(title string) *Movie {
	return &Movie{Title: title}
}

func GetMovieById(id uint, db *gorm.DB, moviesCache *lru_cache.LruCache[uint, Movie]) (Movie, bool) {
	movieFromCache, exists := moviesCache.Get(id)

	if exists == true {
		return movieFromCache, true
	}

	var movieFromDb Movie

	result := db.First(&movieFromDb, id)

	if result.Error == nil {
		moviesCache.Update(id, movieFromDb)

		return movieFromCache, true
	}

	var r Movie

	return r, false
}

func StringToUint(str string) (uint, bool) {
	preId, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		var r uint
		return r, false
	}

	id := uint(preId)

	return id, true
}
