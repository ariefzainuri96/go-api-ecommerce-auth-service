package main

import (
	"log"
	"os"

	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/db"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/logger"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/server"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	gormDb, errGorm := db.NewGorm(os.Getenv("DB_ADDR"))

	if errGorm != nil {
		log.Fatal("Error connecting to gorm database")
		panic("Error connecting to gorm database")
	}

	db, err := db.New(os.Getenv("DB_ADDR"), 30, 30, "10m")

	if err != nil {
		log.Fatal("Error connecting to database")
		panic("Error connecting to database")
	}

	defer db.Close()

	s := server.NewGRPCServer(50051)

	store := store.NewStorage(db, gormDb)

	logger := logger.NewLogger()
	defer logger.Sync()

	if err := s.Start(store, logger); err != nil {
		log.Fatal(err)
	}
}
