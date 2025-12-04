package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() {
	dbURL := os.Getenv("DB_ADDR")
	if dbURL == "" {
		log.Println("⚠️  No DB_ADDR provided, skip migrations")
		return
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Cannot connect to DB: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Cannot create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("❌ Cannot run migrations: %v", err)
	}

	// Try to run all migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Migration error: %v", err)
	}

	log.Println("✅ Migrations applied successfully")
}
