package db

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	_ "github.com/lib/pq"
	"os"
)

func BuildDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
}

func NewClient() *generated.Client {
	client, err := generated.Open("postgres", BuildDSN())
	if err != nil {
		log.Fatalf("❌ failed to open DB: %v", err)
	}
	return client
}
