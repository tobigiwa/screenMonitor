package main

import (
	"context"
	"embed"
	"log"
	"log/slog"
	"strings"

	"agent"
	helperFuncs "pkg/helper"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/*
var assets embed.FS

func main() {

	// logging
	logger, logFile, err := helperFuncs.Logger("desktop.log")
	if err != nil {
		log.Fatalln(err) // exit
	}
	defer logFile.Close()

	slog.SetDefault(logger)

	desktopAgent, err := agent.DesktopAgent(logger)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			log.Fatalln("daemon service is not running", err)
		}
		log.Fatalln("error creating app:", err)
	}

	_, err = desktopAgent.CheckDaemonService()
	if err != nil {
		log.Fatalln("error connecting to daemon service:", err)
	}

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	if err = wails.Run(&options.App{
		Title:  "LiScreMon",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: desktopAgent.Routes(),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	}); err != nil {
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
