package daemon

import (
	db "LiScreMon/cli/daemon/internal/database"
	monitoring "LiScreMon/cli/daemon/internal/monitoring/linux"
	"LiScreMon/cli/daemon/internal/service"
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	badgerDB, err := db.NewBadgerDb(configDir + "/badgerDB/")
	if err != nil {
		log.Fatal(err) // exit
	}

	monitor := monitoring.InitMonitoring(badgerDB)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go service.StartService(socketDir, badgerDB)

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTimer(time.Duration(1) * time.Minute)

	go monitor.WindowChangeTimerFunc(ctx, timer)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	go xevent.Main(monitor.X11Connection) // Start the x11 event loop.
	<-sig
	close(sig)

	// err := monitor.Db.UpdateOpertionOnBuCKET("app", db.ExampleOf_opsFunc)
	// if err != nil {
	// 	fmt.Println("opt failed", err)
	// }

	xevent.Quit(monitor.X11Connection)       // this should always comes first
	cancel()                                 // a different goroutine for managing backing up app usage every minute, fired from monitor
	monitor.CloseWindowChangeCh()            // a different goroutine,closes a channel, this should be after calling the CancelFunc passed to monitor.WindowChangeTimerFunc
	service.ServiceInstance.StopTaskManger() // a different goroutine for managing taskManager, fired from service
	service.SocketConn.Close()
	monitor.Db.Close()

	os.Exit(0)
}
