package daemon

import (
	db "LiScreMon/cli/daemon/internal/database"
	monitoring "LiScreMon/cli/daemon/internal/monitoring/linux"
	"LiScreMon/cli/daemon/internal/service"

	"context"
	"fmt"
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
		log.Fatalln(err) // exit
	}

	// logging
	logger, logFile, err := helperFuncs.Logger(fmt.Sprintf("%s/daemon.log", configDir))
	if err != nil {
		log.Fatalln(err) // exit
	}

	slog.SetDefault(logger)

	// database
	badgerDB, err := db.NewBadgerDb(configDir + "/badgerDB/")
	if err != nil {
		log.Fatalln(err) // exit
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// service
	service, err := service.NewService(badgerDB)
	if err != nil {
		log.Fatalln(err) // exit
	}

	go func() {
		if err := service.StartService(fmt.Sprintf("%s/socket/", configDir), badgerDB); err != nil {
			fmt.Println("error starting service", err)
			sig <- syscall.SIGTERM //if service.StartService fails, send a signal to close the program
		}
	}()

	monitor, err := monitoring.InitMonitoring(badgerDB)
	if err != nil {
		log.Fatalln(err) // exit
	}

	ctx, ctxCancel := context.WithCancel(context.Background())
	timer := time.NewTimer(time.Duration(58) * time.Second)

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
	ctxCancel()                        // a different goroutine for managing backing up app usage every minute, fired from monitor
	monitor.CloseWindowChangeCh()      // a different goroutine,closes a channel, this should be after calling the CancelFunc passed to monitor.WindowChangeTimerFunc

	if !timer.Stop() {
		<-timer.C
	}

	service.StopTaskManger() // a different goroutine for managing taskManager, fired from service
	badgerDB.Close()
	logFile.Close()

	os.Exit(0)
}
