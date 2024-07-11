package agent

import (
	"agent/internal/backend"
	"agent/internal/frontend/desktop"
	"embed"
	"fmt"
	"log/slog"
)

//go:embed internal/frontend/*
var Assets embed.FS

func NewDeskTopApp(logger *slog.Logger) (*backend.App, error) {
	if err := desktop.CreateIndexHTML(); err != nil {
		return nil, fmt.Errorf("error creating index.html needed for desktopApp: %v", err)
	}
	return backend.NewApp(logger)
}
