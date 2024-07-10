package agent

import (
	"agent/internal/backend"
	"log/slog"
)

func NewApp(logger *slog.Logger) (*backend.App, error) {
	return backend.NewApp(logger)
}
