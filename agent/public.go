package agent

import (
	"agent/internal/backend"
	"agent/internal/frontend/desktop"
	"embed"
	"log/slog"
)

//go:embed internal/frontend/*
var Assets embed.FS

func NewDeskTopApp(logger *slog.Logger) (*backend.App, error) {
	if err := desktop.CreateIndexHTML(); err != nil {
		return nil, err
	}
	return backend.NewApp(logger)
}
