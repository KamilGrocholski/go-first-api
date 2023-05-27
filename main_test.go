package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFakeDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func TestingE2E(t *testing.T) {
	db := setupFakeDatabase()
	Start(db)

	t.Run("it should create a movie ang get the movie", func(t *testing.T) {
		router := setupRouter()

		w := httptest.NewRecorder()

		body := struct {
			title string
		}{
			title: "Some Title",
		}

		out, err := json.Marshal(body)
		if err != nil {
			log.Fatal("Marshal error")
		}

		createRequest, _ := http.NewRequest("POST", "/movie/create", bytes.NewBuffer(out))
		router.ServeHTTP(w, createRequest)

		assert.Equal(t, 200, w.Code)

		getRequest, _ := http.NewRequest("GET", "/movie?id=1", nil)
		router.ServeHTTP(w, getRequest)

		expectedGetBody := `{"title": "Some Title"}`

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, expectedGetBody, w.Body.String())
	})
}
