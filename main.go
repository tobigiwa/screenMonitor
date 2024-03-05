package main

import (
	"LiScreMon/store"
	"log"
	"os"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

var (
	X   *xgbutil.XUtil
	err error

	windows []xproto.Window

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

	if err = xproto.ChangeWindowAttributesChecked(X.Conn(), X.RootWin(), xproto.CwEventMask,
		[]uint32{
			xproto.EventMaskFocusChange |
				// xproto.EventMaskPropertyChange |
				// xproto.EventMaskLeaveWindow |
				// xproto.EventMaskEnterWindow |
				xproto.EventMaskSubstructureNotify}).Check(); err != nil {
		log.Fatal("Failed to select notify events for root:", err)
	}

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		rootMapNotifyHandler(X, ev)
	}).Connect(X, X.RootWin())

	// xevent.PropertyNotifyFun(
	// 	func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
	// 		rootPropertyNotifyHandler(X, ev)
	// 	}).Connect(X, X.RootWin())

	if windows, err = currentlyOpenedWindows(X); err != nil {
		log.Fatal(err)
	}

	log.Println("len(windows)=====>>>>>:", len(windows))
	for _, window := range windows {
		name, err := getApplicationName(X, window)
		if err != nil {
			log.Printf("getApplicationName error on window %d:%v\n\n", window, err)
		} else {
			log.Println(window, "====>", name)
		}

		updateWindowInfo(window, name)
	}

	if err := SetDefaultFocusWindow(X); err != nil {
		log.Fatal(err)
	}

	InitMonitoringEvent(X, windows)

	xevent.Main(X)
}
