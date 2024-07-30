package utils

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func Logger(logFileName string, mode bool) (*slog.Logger, *os.File, error) {

	configDir := APP_CONFIG_DIR

	logFile, err := os.Create(filepath.Join(configDir, "logs", logFileName))
	if err != nil {
		return nil, nil, err
	}

	opts := slog.HandlerOptions{
		AddSource: true,
	}

	if mode {
		jsonLogger := slog.NewTextHandler(io.MultiWriter(logFile, os.Stdout), &opts)
		return slog.New(jsonLogger), logFile, nil
	}

	jsonLogger := slog.NewTextHandler(logFile, &opts)
	return slog.New(jsonLogger), logFile, nil
}
