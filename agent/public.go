package agent

import (
	"agent/internal/backend"
	frontend "agent/internal/frontend/components"
	"log/slog"

	"github.com/a-h/templ"
)

func BrowserAgent(logger *slog.Logger) (*backend.App, error) {
	return backend.NewApp(logger)
}

func DesktopAgent(logger *slog.Logger) (*backend.App, error) {
	return backend.NewApp(logger)
}

func IndexPage() templ.Component { // for the desktopApp
	return frontend.IndexPage()
}
