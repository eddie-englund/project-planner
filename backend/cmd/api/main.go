package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/eddie-englund/project-planner/backend/internal/api"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Starting application")

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Error("failed to load .env file", "error", err)
		os.Exit(1)
	}

	router := api.NewRouter(logger)
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
