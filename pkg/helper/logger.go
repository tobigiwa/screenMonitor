package helper

import (
	"io"
	"log/slog"
	"os"
)

func Logger(logFilePath string) (*slog.Logger, *os.File, error) {

	logFile, err := os.Create(logFilePath)
	if err != nil {
		return nil, nil, err
	}

	opts := slog.HandlerOptions{
		AddSource: true,
	}

	jsonLogger := slog.NewTextHandler(io.MultiWriter(logFile, os.Stdout), &opts)

	return slog.New(jsonLogger), logFile, nil
}
