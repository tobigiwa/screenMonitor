package daemon

import (
	db "LiScreMon/cli/daemon/internal/database"
	monitoring "LiScreMon/cli/daemon/internal/monitoring/linux"
	"LiScreMon/cli/daemon/internal/service"

	// "LiScreMon/cli/daemon/internal/service"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	helperFuncs "pkg/helper"

	"github.com/BurntSushi/xgbutil/xevent"
)

func DaemonServiceLinux() {

	// config directory
	configDir, err := helperFuncs.ConfigDir()
	if err != nil {
		log.Fatal(err) // exit
	}

	socketDir := fmt.Sprintf("%s/socket/", configDir)
	logFilePath := fmt.Sprintf("%s/log.txt", configDir)

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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// service
	service, err := service.NewService(badgerDB)
	if err != nil {
		log.Fatal(err) // exit
	}

	go func() {
		if err := service.StartService(socketDir, badgerDB); err != nil {
			time.Sleep(2 * time.Second)
			sig <- syscall.SIGTERM //if service.StartService fails, send a signal to close the program
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.NewTimer(time.Duration(58) * time.Second)

	monitor := monitoring.InitMonitoring(badgerDB)

	go func() {
		monitor.WindowChangeTimerFunc(ctx, timer)
	}()

	go func() {
		xevent.Main(monitor.X11Connection) // Start the x11 event loop.
	}()

	<-sig
	close(sig)

	// err = monitor.Db.UpdateAppInfoManually([]byte("app:Google-chrome"), db.ExampleOf_opsFunc)
	// if err != nil {
	// 	fmt.Println("opt failed", err)
	// }

	xevent.Quit(monitor.X11Connection) // this should always comes first
	cancel()                           // a different goroutine for managing backing up app usage every minute, fired from monitor
	monitor.CloseWindowChangeCh()      // a different goroutine,closes a channel, this should be after calling the CancelFunc passed to monitor.WindowChangeTimerFunc

	if !timer.Stop() {
		<-timer.C
	}

	service.StopTaskManger() // a different goroutine for managing taskManager, fired from service
	badgerDB.Close()

	os.Exit(0)
}
