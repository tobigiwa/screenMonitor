package main

import (
	"LiScreMon/store"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

var (
	X   *xgbutil.XUtil
	err error

	windows = make([]xproto.Window, 0, 10)

	app *X11
)

func main() {

	if X, err = xgbutil.NewConn(); err != nil {
		log.Fatal(err)
	}

	workdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	db, err := store.NewBadgerDb(workdir + "/store/badgerDB/")
	if err != nil {
		log.Fatal(err)
	}

	app = &X11{
		db: db,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigs
		app.db.ReadAll()
		xevent.Quit(X)
		os.Exit(0)
	}()

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

	fmt.Println()

	if err := InitNetActiveWindow(X); err != nil {
		log.Fatal("cannot get InitACtive window", err)
	}

	// Start the event loop.
	xevent.Main(X)
}
