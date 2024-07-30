package monitoring

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

func registerWindow(windowId xproto.Window) {
	xevent.DestroyNotifyFun(destroyNotifyEventFuncWindow).Connect(x11Conn, windowId)
}

func destroyNotifyEventFuncWindow(x11Conn *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
	log.Printf("DESTROY--WINDOW<========Window %d:%s WAS DESTROYED!!! ev.Event:%v========>\n", ev.Window, curSessionNamedWindow[ev.Window], ev.Event)
	xevent.Detach(x11Conn, ev.Window)
}
