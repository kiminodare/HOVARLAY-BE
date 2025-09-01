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
	_ = godotenv.Load()

	// From env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	ssl := os.Getenv("DB_SSLMODE")

	dsn := buildPostgresDSN(host, port, user, pass, name, ssl)

	client, err := generated.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatalf("❌ failed opening connection: %v", err)
	}
	defer func(client *generated.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("❌ failed closing connection: %v", err)
		}
	}(client)

	ctx := context.Background()

	// Migration options
	var opts []schema.MigrateOption
	if os.Getenv("APP_ENV") == "development" {
		opts = append(opts,
			schema.WithDropIndex(true),
			schema.WithDropColumn(true),
		)
	}

	if err := client.Schema.Create(ctx, opts...); err != nil {
		log.Fatalf("❌ migration failed: %v", err)
	}

	log.Println("✅ Migration successful")
}
