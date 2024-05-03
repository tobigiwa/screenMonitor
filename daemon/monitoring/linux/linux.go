package monitoring

import (
	"fmt"
	"log"
	"time"

	"LiScreMon/daemon/repository"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

var (
	// curSessionNamedWindow is a map of all current session "named" windows.
	// An X session is typically a time between login and logout (or restart/shutdown).
	// Only windows with knowm WM_CLASS are added to this map. The X_ID are always unique
	// for a particular window in each session.
	curSessionNamedWindow = make(map[xproto.Window]string, 20)
	netActiveWindowAtom   xproto.Atom
	netClientStackingAtom xproto.Atom
	netActiveWindow       = &netActiveWindowInfo{}
	monitor               X11Monitor
	x11Conn               *xgbutil.XUtil
)

func InitMonitoring(configDir string) X11Monitor {

	// database
	db, err := repository.NewBadgerDb(configDir + "/badgerDB/")
	if err != nil {
		log.Fatal(err) // exit
	}

	// X server connection
	for {
		if x11Conn, err = xgbutil.NewConn(); err != nil { // we wait till we connect to X server
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	monitor = X11Monitor{
		X11Connection: x11Conn,
		Db:            db,
	}

	setRootEventMask(x11Conn)

	registerRootWindowForEvents(x11Conn)

	windows, err := currentlyOpenedWindows(x11Conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("len(windows)=====>>>>>:", len(windows))
	for _, window := range windows {
		name, err := getWindowClassName(x11Conn, window)
		if err != nil {
			continue
		}
		fmt.Println(window, "===========>", name)

		registerWindowForEvents(window)
		addWindowTocurSessionNamedWindowMap(window, name)
	}

	netActiveWindowAtom, netClientStackingAtom = neededAtom()[0], neededAtom()[1]
	netActiveWindow.WindowID = xevent.NoWindow

	return monitor
}
