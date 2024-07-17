package utils

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func Logger(logFileName string) (*slog.Logger, *os.File, error) {

	configDir, err := ConfigDir()
	if err != nil {
		return nil, nil, err
	}

	logFile, err := os.Create(filepath.Join(configDir, "logs", logFileName))
	if err != nil {
		return nil, nil, err
	}

	opts := slog.HandlerOptions{
		AddSource: true,
	}

	jsonLogger := slog.NewTextHandler(io.MultiWriter(logFile, os.Stdout), &opts)
	return slog.New(jsonLogger), logFile, nil
}
