package daemon

import (
	"path/filepath"
	db "smDaemon/daemon/internal/database"
	monitoring "smDaemon/daemon/internal/screen/linux"
	"smDaemon/daemon/internal/service"

	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"utils"

	"github.com/BurntSushi/xgbutil/xevent"
)

func DaemonServiceLinux(logger *slog.Logger) {

	// config directory
	configDir := utils.APP_CONFIG_DIR

	// database
	badgerDB, err := db.NewBadgerDb(filepath.Join(configDir, "badgerDB"))
	if err != nil {
		log.Fatalln(err) // exit
	}

	sig := make(chan os.Signal, 3)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// service
	service, err := service.NewService(badgerDB)
	if err != nil {
		log.Fatalln(err) // exit
	}

	go func() {
		if err := service.StartService(filepath.Join(configDir, "socket"), badgerDB); err != nil {
			log.Println("error starting service", err)
			sig <- syscall.SIGTERM // if service.StartService fails, send a signal to close the program
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
		log.Println("error starting x11 event loop", err)
		sig <- syscall.SIGTERM // if the event loop cannot be started, send a signal to close the program
	}()

	<-sig // awaiting only the first signal

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
	close(sig)
}
