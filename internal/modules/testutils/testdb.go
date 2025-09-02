package testutils

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"testing"

	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
	_ "github.com/lib/pq"
)

func NewTestDB(t *testing.T) *generated.Client {
	t.Helper()

	err := godotenv.Load(filepath.Join("E:\\test\\HOVARLAY\\BE", ".env.local"))
	if err != nil {
		t.Fatalf("failed loading .env file: %v", err)
	}

	fmt.Println(os.Getenv("DB_HOST"))

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	ssl := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, ssl,
	)

	client, err := generated.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed opening postgres db: %v", err)
	}

	// Reset schema supaya test bersih
	if err := client.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed creating schema: %v", err)
	}

	return client
}
