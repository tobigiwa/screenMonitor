package main

import (
	"LiScreMon/store"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
)

func (x11 *X11) rootMapNotifyHandler(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
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
		name = "name-not-found"
	}

	if name != "name-not-found" {
		log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionOpenedWindowMap(ev.Window, name)
		addWindowTocurSessionNamedWindowMap(ev.Window, name)
	}
}

func (x11 *X11) rootPropertyNotifyHandler(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent, netActiveWindowAtom xproto.Atom) {

	if ev.Atom == netActiveWindowAtom {
		if activeWin, err := ewmh.ActiveWindowGet(X); (err == nil) && (activeWin != 0) { // 0 is root, to much noise

			if activeWin != netActiveWindow.WindowID { // a window has become active
				s := store.ScreenTime{
					AppName: netActiveWindow.WindowName,
					Type:    store.Active,
					Time:    time.Since(netActiveWindow.TimeStamp).Hours(),
				}

				fmt.Printf("New active window ID =====> %v:%v:%v\n", activeWin, curSessionNamedWindow[activeWin], curSessionOpenedWindow[activeWin].Name)
				fmt.Printf("time elapsed for window %v:%v was %v\n\n", activeWin, curSessionNamedWindow[activeWin], time.Since(netActiveWindow.TimeStamp).Hours())

				if _, exists := curSessionNamedWindow[activeWin]; !exists { // if name does not already exist in curSessionNamedWindow, include it.
					fmt.Println("it got to resolve name for window")
					if name, err := getWindowClassName(X, activeWin); err == nil {
						curSessionNamedWindow[activeWin] = name
						fmt.Printf("window:%v ====> name:%v now added in curSessionNamedWindow\n", activeWin, name)
					}
				}

				netActiveWindow.WindowID = activeWin
				netActiveWindow.WindowName = curSessionNamedWindow[activeWin]
				netActiveWindow.TimeStamp = time.Now()

				if err := x11.db.WriteUsage(s); err != nil {
					log.Fatalf("focusOutEventHandler:write to db error:%v", err)
				}

			}
		}

	}
}
