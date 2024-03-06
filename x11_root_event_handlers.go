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
		ok   bool
	)

	fmt.Printf("\nrootMapNotifyHandler ev.window:%v ======++++++====> ev.event:%v\n", ev.Window, ev.Event)

	if name, ok = curSessionNamedWindow[ev.Window]; ok {
		fmt.Printf("window:%v name resolved from window ITSELF %s:%v\n", ev.Window, name, ev.Event)
		goto jump
	}

	fmt.Println("--->--->--->")

	if name, err = getApplicationName(X, ev.Window); err != nil && name == "" {

		fmt.Printf("getApplicationName:error on window %d:\n %v\n", ev.Window, err)

		if name, err = checkQueryTreeForParent(X, ev.Window); err != nil {

			fmt.Printf("checkQueryTreeForParent:error on window %v: error: %v\n", ev.Window, err)

			name = "name-not-found"

			list, err := currentlyOpenedWindows(X)
			if err != nil {
				log.Fatalf("err in getting all windows in rootMapNotifyHandler %v\n", err)
			}

			for _, window := range list {
				if window == ev.Window {
					name, err := getApplicationName(X, ev.Window)
					if err != nil {
						log.Println(err)
					}

					fmt.Println("WE GOT A NAME======================>,", name)

				}
			}
		}
	}

jump:
	if _, exists := curSessionNamedWindow[ev.Window]; !exists && (name != "name-not-found") {
		curSessionNamedWindow[ev.Window] = name
	}
	log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)

	addWindowTocurSessionOpenedWindowMap(ev.Window, name)
}

func rootPropertyNotifyHandler(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent, netActiveWindowAtom xproto.Atom) {

	// fmt.Printf("\nrootPropertyNotifyHandler =====>widow:%v... got atom %v, expecting atom %v\n\n", ev.Window, ev.Atom, netActiveWindowAtom)

	if ev.Atom == netActiveWindowAtom {
		fmt.Println("Active window changed.")
		activeWin, err := ewmh.ActiveWindowGet(X)
		if err != nil {
			log.Printf("Failed to get active window ID: %v", err)
		} else {
			fmt.Printf("New active window ID =====> %v:%v:%v\n", activeWin, curSessionNamedWindow[activeWin], curSessionOpenedWindow[activeWin].Name)

			if activeWin == globalFocusEvent.WindowID {
				fmt.Printf("###########ACTIVE:%v:%v:%v ===== FOCUS:%v:%v###########\n", activeWin, curSessionNamedWindow[activeWin], curSessionOpenedWindow[activeWin].Name, globalFocusEvent.WindowID, globalFocusEvent.AppName)
			} else {
				fmt.Printf("$$$$$$$$$$$$ACTIVE:%v:%v:%v !=!=!=!=!=!=!= FOCUS:%v:%v$$$$$$$$$$$$\n", activeWin, curSessionNamedWindow[activeWin], curSessionOpenedWindow[activeWin].Name, globalFocusEvent.WindowID, globalFocusEvent.AppName)
			}

			if _, exists := curSessionNamedWindow[activeWin]; !exists {
				name, err := getApplicationName(X, activeWin)
				if err != nil {
					log.Printf("getApplicationName error on window %d:%v\n\n", activeWin, err)
				} else {
					curSessionNamedWindow[activeWin] = name
					log.Printf("%v ====> %v and now added in curSessionNamedWindow", activeWin, name)
				}
			}
		}
	}
}
