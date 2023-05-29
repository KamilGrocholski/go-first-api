package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikalsqwe/go-api/config/env"
	"github.com/mikalsqwe/go-api/models"
	"github.com/mikalsqwe/go-api/utils"
)

const BASE_URL string = "http://localhost:3000"

func ComposeUrl(v string) string {
	return BASE_URL + v
}

func TestBooksRouter(t *testing.T) {
	router := SetupServer(env.TEST_DATABASE_URL)

	t.Run("it should create a book", func(t *testing.T) {
		book := models.CreateBookPayload{
			Title: "Nowy tytuł",
		}

		bodyData, _ := json.Marshal(book)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, ComposeUrl("/books/create"), bytes.NewReader(bodyData))
		router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Error("not 200")
		}
	})

	t.Run("it should get the created book", func(t *testing.T) {
		payload := models.GetBookByIdPayload{
			Id: 1,
		}

		w := httptest.NewRecorder()

		path := fmt.Sprint("/books/", payload.Id)

		req, _ := http.NewRequest(http.MethodGet, ComposeUrl(path), nil)
		router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Error("not 200")
		}
	})

	t.Run("it should get a page of books", func(t *testing.T) {
		payload := utils.PaginationPayload{
			Limit:  10,
			Page:   1,
			SortBy: "id",
		}

		w := httptest.NewRecorder()

		path := fmt.Sprintf("/books?limit=%v&page=%v&sort_by=%v", payload.Limit, payload.Page, payload.SortBy)

		req, _ := http.NewRequest(http.MethodGet, ComposeUrl(path), nil)
		router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Error("not 200")
		}
	})

	t.Run("it should update a book", func(t *testing.T) {
		payload := models.UpdateBookPayload{
			Id:    1,
			Title: "Nowy tytuł",
		}

		w := httptest.NewRecorder()

		bodyData, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPatch, ComposeUrl("/books/update"), bytes.NewReader(bodyData))
		router.ServeHTTP(w, req)

		if http.StatusOK != w.Code {
			t.Error("not 200")
		}
	})
}
