package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/kiminodare/HOVARLAY-BE/ent/generated"
)

func buildPostgresDSN(host, port, user, pass, name, ssl string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, name, ssl,
	)
}

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using environment variables")
	}

	dsn := buildPostgresDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	client, err := generated.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatalf("❌ failed opening connection: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	var opts []schema.MigrateOption

	appEnv := os.Getenv("APP_ENV")

	switch appEnv {
	case "local":
		opts = append(opts,
			schema.WithDropColumn(true),
			schema.WithDropIndex(true),
			schema.WithForeignKeys(true),
		)
	case "development", "production":
		opts = append(opts,
			schema.WithForeignKeys(true),
		)
	default:
		log.Fatalf("❌ Unknown APP_ENV: %s", appEnv)
	}

	if err := client.Schema.Create(ctx, opts...); err != nil {
		log.Fatalf("❌ migration failed: %v", err)
	}

	log.Printf("✅ Migration successful for %s environment", appEnv)
}
