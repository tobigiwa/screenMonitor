package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

func registerWindowForEvents(windowId xproto.Window) {

	xevent.DestroyNotifyFun(func(xu *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		if window, ok := curSessionOpenedWindow[ev.Window]; ok {
			deleteWindowFromcurSessionOpenedWindowMap(ev.Window)
			log.Printf("<========Window %d:%s was destroyed ev.Event:%v========>\n", ev.Window, window.Name, ev.Event)
		}
	}).Connect(X, windowId)

	xevent.FocusInFun(func(X *xgbutil.XUtil, ev xevent.FocusInEvent) {
		fmt.Printf("+++++++Focus in on window %d:%+v\n", ev.Event, curSessionOpenedWindow[ev.Event].Name)
		app.focusInEventHandler(ev)
	}).Connect(X, windowId)

	xevent.FocusOutFun(func(X *xgbutil.XUtil, ev xevent.FocusOutEvent) {
		fmt.Printf("----Focus out on window %d:%+v", ev.Event, curSessionOpenedWindow[ev.Event].Name)
		app.focusOutEventHandler(ev)
	}).Connect(X, windowId)

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		fmt.Printf("%v\n", strings.Repeat("*", 15))
		fmt.Printf("window in the curSessionNamedWindow:%v\n", curSessionNamedWindow[ev.Window])
	}).Connect(X, windowId)

	// xevent.UnmapNotifyFun(func(X *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
	// 	window, ok := curSessionOpenedWindow[ev.Window]
	// 	if !ok {
	// 		fmt.Printf("\n*****name should have been in the Map, unmap*****\n\n")
	// 		return
	// 	}
	// 	log.Printf("Window %d ===> %s was unmapped \n", ev.Window, window.Name)
	// }).Connect(X, windowId)

}
