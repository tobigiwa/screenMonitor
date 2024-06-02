package daemon

import (
	db "LiScreMon/cli/daemon/internal/database"
	monitoring "LiScreMon/cli/daemon/internal/monitoring/linux"
	"LiScreMon/cli/daemon/internal/service"
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
	socketDir := configDir + "/socket/"
	logFilePath := configDir + "/log.txt"

	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal(err) // exit
	}

	// logging
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	db, err := db.NewBadgerDb(configDir + "/badgerDB/")
	if err != nil {
		log.Fatal(err) // exit
	}

	monitor := monitoring.InitMonitoring(db)

	signal1 := make(chan os.Signal, 1)
	signal.Notify(signal1, os.Interrupt, syscall.SIGTERM)

	go service.StartService(socketDir, db)

	go func() {
		<-signal1
		close(signal1)

		// for _, v := range [2]string{"day", "week"} {
		// 	fmt.Println("at", v)
		// 	monitor.Db.DeleteBucket(v)
		// }

		// monitor.Db.UpdateOpertionOnBuCKET("app", database.DeleteCmdOnAppInfo)
		xevent.Quit(monitor.X11Connection)
		service.SocketConn.Close()
		monitor.Db.Close()

		os.Exit(0)
	}()

	// Start the event loop.
	xevent.Main(monitor.X11Connection)
}
