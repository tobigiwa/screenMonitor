package agent

import (
	"agent/internal/backend"
	frontend "agent/internal/frontend/components"
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
)

//go:embed internal/frontend/*
var FrontendDirForWailsDev embed.FS

func BrowserAgent(logger *slog.Logger) (*backend.App, error) {
	return backend.NewApp(logger)
}

func DesktopAgent(logger *slog.Logger) (*backend.App, error) {
	if err := createIndexHTML(); err != nil {
		return nil, fmt.Errorf("error creating index.html needed for desktopApp: %v", err)
	}
	return backend.NewApp(logger)
}

func createIndexHTML() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	p := fmt.Sprintf("%s/frontend/index.html", cwd)
	indexHTML, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}

	if err = frontend.IndexPage().Render(context.Background(), indexHTML); err != nil {
		return fmt.Errorf("could not generate html: %v", err)
	}

	return nil
}
