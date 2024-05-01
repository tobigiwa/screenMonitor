package daemon

import (
	monitoring "LiScreMon/daemon/monitoring/linux"
	"LiScreMon/daemon/service"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/xgbutil/xevent"
)

func DaemonService() {

	// config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err) // exit
	}

	configDir := homeDir + "/liScreMon"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal(err) // exit
	}

	// logging
	logFile, err := os.OpenFile(configDir+"/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err) // exit
	}
	defer logFile.Close()

	opts := slog.HandlerOptions{
		AddSource: true,
	}
	jsonLogger := slog.NewJSONHandler(io.MultiWriter(logFile, os.Stdout), &opts)
	logger := slog.New(jsonLogger)
	slog.SetDefault(logger)

	monitor := monitoring.InitMonitoring(configDir)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func(signal chan os.Signal) {
		<-signal
		monitor.Db.Close()
		xevent.Quit(monitor.X11Connection)
		os.Exit(0)
	}(sigs)

	go service.StartService(homeDir, sigs)

	log.Println("LiScreMon started...")
	// Start the event loop.
	xevent.Main(monitor.X11Connection)
}
