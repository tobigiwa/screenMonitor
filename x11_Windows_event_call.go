package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

func registerWindowForEvents(windowId xproto.Window) {

	xevent.DestroyNotifyFun(func(xu *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		if window, ok := curSessionOpenedWindow[ev.Window]; ok {
			deleteWindowFromcurSessionOpenedWindowMap(ev.Window)
			log.Printf("WINDOW<========Window %d:%s was destroyed ev.Event:%v========>\n", ev.Window, window.Name, ev.Event)
		}
		xevent.Detach(X, windowId)
	}).Connect(X, windowId)

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		fmt.Printf("\nWINDOWMapNotifyHandler ev.window:%v ======++++++====> ev.event:%v\n", ev.Window, ev.Event)
		app.windowMapNotifyHandler(X, ev)
	}).Connect(X, windowId)

	xevent.VisibilityNotifyFun(func(xu *xgbutil.XUtil, ev xevent.VisibilityNotifyEvent) {
		fmt.Printf("\nWINDOW::::window:%v:name-%v got visibility of state--- %v\n", ev.Window, curSessionNamedWindow[ev.Window], ev.State)
	}).Connect(X, windowId)

}
