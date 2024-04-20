package main

import (
	"LiScreMon/store"
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
	X   *xgbutil.XUtil
	err error

	windows = make([]xproto.Window, 0, 10)

	netActiveWindowAtom   xproto.Atom
	netClientStackingAtom xproto.Atom

	app *X11
)

func main() {

	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err) // fail
	}

	dirPath := path + "/liScreMon"
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatal(err) // fail
	}

	logFile, err := os.OpenFile(dirPath+"/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err) // fail
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
		os.Exit(0)
	}()

	netActiveWindowAtom, netClientStackingAtom = neededAtom()[0], neededAtom()[1]

	setRootEventMask(X)

	registerRootWindowForEvent(X)

	if windows, err = currentlyOpenedWindows(X); err != nil {
		log.Fatal(err)
	}

	log.Println("len(windows)=====>>>>>:", len(windows))
	for _, window := range windows {
		name, err := getWindowClassName(X, window)
		if err != nil {
			log.Printf("getWindowClassName error on window %d:%v\n", window, err)
			continue
		}

		log.Println(window, "===========>", name)

		addWindowTocurSessionOpenedWindowMap(window, name)
		addWindowTocurSessionNamedWindowMap(window, name)
	}

	netActiveWindow.WindowID = xevent.NoWindow

	// Start the event loop.
	xevent.Main(X)
}
