package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mikalsqwe/go-api/book"
)

func TestApp(t *testing.T) {
	t.Run("it should do smth", func(t *testing.T) {
		t.Log("end")
		Start("test.db")

		bookPayload := &book.CreateBookPayload{
			Title:       "Some title",
			Description: "Some description",
			AuthorsIds:  make([]string, 0),
		}

		out, stringifyErr := json.Marshal(bookPayload)

		if stringifyErr != nil {
			panic("")
		}

		res, err := http.Post("/create/book", "application/json", bytes.NewBuffer(out))
		if err != nil {
			t.Errorf("error")
		}

		res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("not 200")
		}
		t.Log("end")
	})
}
