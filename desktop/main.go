package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"

	desktopApp "agent"
	helperFuncs "pkg/helper"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	// logging
	logger, logFile, err := helperFuncs.Logger("desktop.log")
	if err != nil {
		log.Fatalln(err) // exit
	}
	defer logFile.Close()

	slog.SetDefault(logger)

	agent, err := desktopApp.NewApp(logger)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Fatalln("daemon service is not running", err)
		}
		log.Fatalln("error creating app:", err)
	}

	_, err = agent.CheckDaemonService()
	if err != nil {
		log.Fatalln("error connecting to daemon service:", err)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "LiScreMon",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  desktopApp.Assets,
			Handler: agent.Routes(),
			// Middleware: func(next http.Handler) http.Handler {
			// 	&http.ServeMux{}.NotFound(next.ServeHTTP)
			// 	return r
			// },
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
