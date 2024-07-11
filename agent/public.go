package agent

import (
	"agent/internal/backend"
	"embed"
	"log/slog"
)

//go:embed internal/frontend/*
var Assets embed.FS

func NewApp(logger *slog.Logger) (*backend.App, error) {
	return backend.NewApp(logger)
}
