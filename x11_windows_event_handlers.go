package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

var globalFocusEvent focusEvent

func (x11 X11) focusInEventHandler(fv xevent.FocusInEvent) {

	if fv.Event != globalFocusEvent.WindowID {
		globalFocusEvent.Time = time.Now()
		globalFocusEvent.WindowID = fv.Event

	}
}

func (x11 X11) focusOutEventHandler(fv xevent.FocusOutEvent) {

	var duration time.Duration
	currentFocusWIndow := globalFocusEvent.WindowID
	if fv.Event == currentFocusWIndow {
		duration = time.Since(globalFocusEvent.Time)

		fmt.Printf("time elapsed for window %v:%v was %v\n", currentFocusWIndow, allCurentlyOpenedWindowMap[currentFocusWIndow].Name, duration.Minutes())

		// s := store.ScreenTime{
		// 	AppName: globalFocusEvent.AppName,
		// 	Type:    store.Active,
		// 	Time:    duration.Minutes(),
		// }

		// if err := x11.db.WriteUsuage(s); err != nil {
		// 	log.Fatalf("focusOutEventHandler:write to db error:%v", err)
		// }
	}
}

func (x11 X11) propertyNotifyEventHandler(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
	atom, err := xprop.Atm(X, "_NET_WM_NAME")
	if err != nil {
		log.Fatalf("Could not get _NET_WM_NAME atom: %v", err)
	}

	atom1, err := xprop.Atm(X, "WM_NAME")
	if err != nil {
		log.Fatalf("Could not get _NET_WM_NAME atom: %v", err)
	}

	atom2, err := xprop.Atm(X, "WM_CLASS")
	if err != nil {
		log.Fatalf("Could not get _NET_WM_NAME atom: %v", err)
	}

	atom3, err := xprop.Atm(X, "_NET_WM_VISIBLE_NAME")
	if err != nil {
		log.Fatalf("Could not get _NET_WM_NAME atom: %v", err)
	}

	if ev.Atom == atom || ev.Atom == atom1 || ev.Atom == atom2 || ev.Atom == atom3 {
		fmt.Println("FOOOOOOOOONAAAAAAALLLLLLLLLLLYYY---WINDOW")
		wmName, err := ewmh.WmNameGet(X, ev.Window)
		if err == nil && wmName != "" {
			fmt.Printf("-+-+-+-+ProrpertyNotify: The WmNameGet on windowID:%v -----> %+v\n", ev.Window, wmName)
			if name, exist := allSessionOpenedAndNamedWindow[ev.Window]; !exist {
				allSessionOpenedAndNamedWindow[ev.Window] = wmName
			} else {
				fmt.Printf("name was previously %s of window:%v and name name is %v\n", name, ev.Window, wmName)
			}
		} else {
			log.Printf("-+-+-+-+ProrpertyNotify: Could not get _NET_WM_NAME for window %d: %v\n", ev.Window, err)
		}
	} else {
		fmt.Printf("In propertyNotify:WINDOW but atom was %v and wanted atom is %v or %v or %v or %v\n", ev.Atom, atom, atom1, atom2, atom3)
	}

}
