package monitoring

import (
	"log"
	"pkg/types"
	"time"

	db "LiScreMon/cli/daemon/internal/database"

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

	fixtyEightSecs = time.Duration(58) * time.Second
)

func InitMonitoring(db *db.BadgerDBStore) X11Monitor {

	var err error
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
		X11Connection:  x11Conn,
		Db:             db,
		windowChangeCh: make(chan types.GenericKeyValue[xproto.Window, float64], 1),
	}

	setRootEventMask(x11Conn)

	registerRootWindowForEvents(x11Conn)

	windows, err := currentlyOpenedWindows(x11Conn)
	if err != nil {
		log.Fatal(err)
	}

	for _, window := range windows {
		getWindowClassName(x11Conn, window)
		registerWindowForEvents(window)
	}

	netActiveWindowAtom, netClientStackingAtom = neededAtom()[0], neededAtom()[1]
	netActiveWindow.WindowID = xevent.NoWindow

	return monitor
}
