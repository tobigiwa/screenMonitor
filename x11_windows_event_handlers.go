package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

func (x11 *X11) windowMapNotifyHandler(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
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

	if name, err = getWindowClassName(X, ev.Window); err != nil && (name != ""){
		fmt.Printf("getWindowClassName:error on window %d:\n %v\n", ev.Window, err)
		name = "name-not-found"
	}

	if name != "name-not-found" {
		log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionOpenedWindowMap(ev.Window, name)
		addWindowTocurSessionNamedWindowMap(ev.Window, name)
	}
}
