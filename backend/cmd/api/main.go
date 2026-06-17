package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/eddie-englund/project-planner/backend/internal/api"
	db "github.com/eddie-englund/project-planner/backend/internal/db/generated"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Starting application")

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Error("failed to load .env file", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := db.New(pool)
	router := api.NewRouter(logger, queries)
	addr := ":8080"
	logger.Info("server listening", "addr", "http://localhost"+addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
