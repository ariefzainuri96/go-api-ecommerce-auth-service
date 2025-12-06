package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

func RunServer(ctx context.Context, cfg Config, logger *zap.Logger) error {
	mux := http.NewServeMux()

	// register health check
	mux.Handle("GET /health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(map[string]string{
			"status": "OK",
		})
		w.Write(jsonResponse)
	}))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.HTTPPort),
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	log.Printf("Server has started on %v", os.Getenv("PORT"))

	// Start server in goroutine so we can watch ctx
	serverErrCh := make(chan error, 1)
	go func() {
		logger.Info("starting http server", zap.Int("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	select {
	case <-ctx.Done():
		// Shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logger.Info("http server shutting down")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("http server shutdown error", zap.Error(err))
			return err
		}
		logger.Info("http server stopped")
		return nil
	case err := <-serverErrCh:
		return err
	}
}
