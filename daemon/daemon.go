package daemon

import (
	"LiScreMon/daemon/store"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

var (
	// curSessionNamedWindow is a map of all current session "named" windows.
	// An X session is typically a time between login and logout (or restart/shutdown).
	// Only windows with knowm WM_CLASS are added to this map. The X_ID are always unique
	// for a particular window in each session.
	curSessionNamedWindow = make(map[xproto.Window]string)
	X                     *xgbutil.XUtil
	netActiveWindowAtom   xproto.Atom
	netClientStackingAtom xproto.Atom
	netActiveWindow       = &netActiveWindowInfo{}
	app                   *X11
)

func DaemonService() {

	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err) // exit
	}

	dirPath := path + "/liScreMon"
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatal(err) // exit
	}

	logFile, err := os.OpenFile(dirPath+"/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err) // exit
	}
	defer logFile.Close()

	// logging
	opts := slog.HandlerOptions{
		AddSource: true,
	}

	jsonLogger := slog.NewJSONHandler(logFile, &opts)
	logger := slog.New(jsonLogger)
	slog.SetDefault(logger)

	// database
	db, err := store.NewBadgerDb(dirPath + "/badgerDB/")
	if err != nil {
		log.Fatal(err) // fail
	}

	// X server connection
	for {
		if X, err = xgbutil.NewConn(); err != nil { // we wait till we connect to X server
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	app = &X11{
		db: db,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		xevent.Quit(X)
		fmt.Println()
		app.db.ReadAll()
		os.Exit(0)
	}()

	netActiveWindowAtom, netClientStackingAtom = neededAtom()[0], neededAtom()[1]

	setRootEventMask(X)

	registerRootWindowForEvents(X)

	windows, err := currentlyOpenedWindows(X)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("len(windows)=====>>>>>:", len(windows))
	for _, window := range windows {
		name, err := getWindowClassName(X, window)
		if err != nil {
			continue
		}
		fmt.Println(window, "===========>", name)

		registerWindowForEvents(window, name)
		addWindowTocurSessionNamedWindowMap(window, name)
	}

	netActiveWindow.WindowID = xevent.NoWindow

	log.Println("LiScreMon started...")

	// Start the event loop.
	xevent.Main(X)
}
