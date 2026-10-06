// Command inventory-transfer-service moves stock of a SKU between warehouse locations.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const defaultDatabaseURL = "postgres://inventory:inventory@localhost:5432/inventory" //nolint:gosec // local development default

func main() {
	if err := run(); err != nil {
		slog.Error("service stopped", "err", err)
		os.Exit(1)
	}
}

func run() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}

	store, err := NewStore(context.Background(), databaseURL)
	if err != nil {
		return err
	}
	defer store.Close()

	h := &Handler{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", h.health)
	mux.HandleFunc("GET /api/inventory/{sku}", h.inventory)
	mux.HandleFunc("POST /api/transfer", h.transfer)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	slog.Info("inventory-transfer-service listening", "addr", srv.Addr)
	return srv.ListenAndServe()
}
