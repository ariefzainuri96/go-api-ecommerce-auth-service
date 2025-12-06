package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/db"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/logger"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/server"
	"github.com/ariefzainuri96/go-api-ecommerce-auth-service/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config holds application configuration (simplified)
type Config struct {
	GRPCPort    int
	HTTPPort    int
	ShutdownTTL time.Duration
}

// loadConfig loads config from environment (simple)
func loadConfig() Config {
	// In real project use env parsing lib (envconfig/viper)
	grpcPort := 50051
	httpPort := 9000
	ttl := 15 * time.Second

	if v := os.Getenv("GRPC_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &grpcPort)
	}
	if v := os.Getenv("HTTP_PORT"); v != "" {
		fmt.Sscanf(v, "%d", &httpPort)
	}
	if v := os.Getenv("SHUTDOWN_TTL"); v != "" {
		var s int
		fmt.Sscanf(v, "%d", &s)
		ttl = time.Duration(s) * time.Second
	}

	return Config{
		GRPCPort:    grpcPort,
		HTTPPort:    httpPort,
		ShutdownTTL: ttl,
	}
}

func main() {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
			return
		}
	}

	cfg := loadConfig()

	db.RunMigrations()

	// get gormdb
	gorm, errGorm := db.NewGorm(os.Getenv("DB_ADDR"))

	if errGorm != nil {
		log.Fatal("Error connecting to gorm database")
		panic("Error connecting to gorm database")
	}

	// get std db
	db, err := db.New(os.Getenv("DB_ADDR"), 30, 30, "10m")

	if err != nil {
		log.Fatal("Error connecting to database")
		panic("Error connecting to database")
	}

	defer db.Close()

	// create grpc server instance
	s := server.NewGRPCServer(cfg.GRPCPort)

	// setup zap logger
	logger := logger.NewLogger()
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// WaitGroup to wait for servers to stop
	var wg sync.WaitGroup

	// run grpc
	wg.Add(1)
	go func() {
		defer wg.Done()

		store := store.NewStorage(db, gorm)

		if err := s.Start(ctx, store, logger); err != nil {
			logger.Error("grpc server stopped with error", zap.Error(err))
			cancel()
		} else {
			logger.Info("grpc server stopped gracefully")
		}
	}()

	// run server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := RunServer(ctx, cfg, logger); err != nil {
			logger.Error("http server stopped with error", zap.Error(err))
			cancel()
		} else {
			logger.Info("http server stopped")
		}
	}()

	// ---------------------------------------------------------
	// Graceful shutdown on OS signals
	// ---------------------------------------------------------
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Info("context canceled, starting shutdown")
	case sig := <-sigCh:
		logger.Info("signal received, starting shutdown", zap.String("signal", sig.String()))
	}

	// start shutdown procedure
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTTL)
	defer shutdownCancel()

	// Let goroutines handle ctx cancellation — we call cancel to notify them
	cancel()

	// Wait for background goroutines to finish or timeout
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		logger.Info("all servers stopped gracefully")
	case <-shutdownCtx.Done():
		logger.Warn("shutdown timed out, forcing exit")
	}

	logger.Info("shutdown complete")
}
