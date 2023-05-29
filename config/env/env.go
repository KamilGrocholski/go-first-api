package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadDotenvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("dotenv load error")
	}
}

func envString(key string) string {
	loadDotenvFile()

	value, exists := os.LookupEnv(key)

	if exists == false {
		log.Fatalf("env: %v is not declared", key)
	}

	if value == "" {
		log.Fatalf("env: %v is empty", key)
	}

	return value
}

var (
	DATABASE_URL      string = envString("DATABASE_URL")
	PORT              string = envString("PORT")
	TEST_DATABASE_URL string = envString("TEST_DATABASE_URL")
)
