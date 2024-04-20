package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

func registerWindow(windowId xproto.Window) {

	xevent.DestroyNotifyFun(destroyNotifyEventFuncWindow).Connect(X, windowId)
}

func destroyNotifyEventFuncWindow(xu *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
	log.Printf("WINDOW<========Window %d:%s WAS DESTROYED!!! ev.Event:%v========>\n", ev.Window, curSessionNamedWindow[ev.Window], ev.Event)
	xevent.Detach(X, ev.Window)
}
