package main

import (
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

func (x11 *X11) windowMapNotifyHandler(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {

	name, _ := getWindowClassName(X, ev.Window)

	if name != "" {
		log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionNamedWindowMap(ev.Window, name)
	}
}
