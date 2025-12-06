package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	gormDb *gorm.DB
}

func (d *GormDB) ExecWithTimeout(ctx context.Context, fn func(tx *gorm.DB) error) error {
    ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
    defer cancel()

    return fn(d.gormDb.WithContext(ctx))
}

func NewGorm(addr string) (*GormDB, error) {
	// Use your existing DSN (Data Source Name) / connection string
	// Example DSN: "host=localhost user=user password=pass dbname=ecommerce-db port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to the database with GORM: %v", err)
		return &GormDB{}, err
	}

	log.Println("Database connection successfully established with GORM.")

	return &GormDB{
		gormDb: db,
	}, nil
}

func New(addr string, maxOpenCons, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)

	if err != nil {
		log.Println("openError:", err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenCons)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)

	if err != nil {
		log.Println("parseDurationError:", err.Error())
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		log.Println("pingError:", err.Error())
		return nil, err
	}

	return db, nil
}
