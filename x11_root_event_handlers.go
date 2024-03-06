package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
)

func rootMapNotifyHandler(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
	var (
		name string
		err  error
	)

	if name, exists := curSessionNamedWindow[ev.Window]; exists {
		fmt.Printf("window:%v name resolved from window ITSELF %s:%v\n", ev.Window, name, ev.Event)
		log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionOpenedWindowMap(ev.Window, name)
		return
	}

	if name, err = getWindowClassName(X, ev.Window); err != nil {
		fmt.Printf("getWindowClassName:error on window %d:\n %v\n", ev.Window, err)
		if name, err = checkQueryTreeForParent(X, ev.Window); err != nil {
			fmt.Printf("checkQueryTreeForParent:error on window %v: error: %v\n", ev.Window, err)
			name = "name-not-found"
		}
	}

	if name != "name-not-found" {
		log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionOpenedWindowMap(ev.Window, name)
		addWindowTocurSessionNamedWindowMap(ev.Window, name)
	}
}

func rootPropertyNotifyHandler(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent, netActiveWindowAtom xproto.Atom) {
	
	if ev.Atom == netActiveWindowAtom {
		if activeWin, err := ewmh.ActiveWindowGet(X); (err == nil) && (activeWin != 0) {
			fmt.Printf("New active window ID =====> %v:%v:%v\n\n", activeWin, curSessionNamedWindow[activeWin], curSessionOpenedWindow[activeWin].Name)
			if _, exists := curSessionNamedWindow[activeWin]; !exists { // if name does not already exist in curSessionNamedWindow, include it.
				if name, err := getWindowClassName(X, activeWin); err == nil {
					curSessionNamedWindow[activeWin] = name
					fmt.Printf("window:%v ====> name:%v now added in curSessionNamedWindow\n", activeWin, name)

				}
			}

		}
	}
}
