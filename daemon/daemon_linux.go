package daemon

import (
	"LiScreMon/daemon/internal/database/repository"
	monitoring "LiScreMon/daemon/internal/monitoring/linux"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/xgbutil/xevent"
)

func DaemonServiceLinux() {

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
	jsonLogger := slog.NewTextHandler(io.MultiWriter(logFile, os.Stdout), &opts)
	logger := slog.New(jsonLogger)
	slog.SetDefault(logger)

	// database
	db, err := repository.NewBadgerDb(configDir + "/badgerDB/")
	if err != nil {
		log.Fatal(err) // exit
	}

	monitor := monitoring.InitMonitoring(db)

	sig1 := make(chan os.Signal, 1)
	signal.Notify(sig1, os.Interrupt, syscall.SIGTERM)

	// ctx, cancel := context.WithCancel(context.Background())
	// go service.StartService(ctx, configDir, db)

	go func() {
		<-sig1
		close(sig1)
		// cancel()
		if err := syscall.Unlink(configDir + "/socket/" + "liScreMon.sock"); err != nil {
			fmt.Println("unlink", err)
		}

		data, err := monitor.Db.GetWeeklyScreenStats(repository.Active, "2024-05-04")
		if err != nil {
			log.Println("error:", err)
		}
		for _, value := range data {
			fmt.Println(value.Key, value.Value)
		}

		monitor.Db.Close()

		xevent.Quit(monitor.X11Connection)
		os.Exit(0)
	}()

	log.Println("LiScreMon started...")
	// Start the event loop.
	xevent.Main(monitor.X11Connection)
}

// sig2 := make(chan os.Signal, 1)
// signal.Notify(sig2, os.Interrupt, syscall.SIGTERM)
