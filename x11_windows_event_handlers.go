package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

func (x11 *X11) windowMapNotifyHandler(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {

	if name, ok := curSessionNamedWindow[ev.Window]; ok {
		fmt.Printf("window:%v name resolved from window ITSELF %s:%v\n", ev.Window, name, ev.Event)
		log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionOpenedWindowMap(ev.Window, name)
		return
	}

	name, _ := getWindowClassName(X, ev.Window)

	if name != "" {
		log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionOpenedWindowMap(ev.Window, name)
		addWindowTocurSessionNamedWindowMap(ev.Window, name)
	}
}
