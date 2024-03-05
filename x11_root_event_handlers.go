package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
)

func rootMapNotifyHandler(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
	var (
		name string
		err  error
		ok   bool
	)

	fmt.Printf("\nrootMapNotifyHandler ev.window:%v ======++++++====> ev.event:%v\n", ev.Window, ev.Event)

	if name, ok = allSessionOpenedAndNamedWindow[ev.Event]; ok {
		fmt.Printf("window:%v name resolved from parent %s:%v\n", ev.Window, name, ev.Event)
		goto jump
	}

	fmt.Println("\nwindow not found in allSessionOpenedAndNamedWindow")

	if name, err = getApplicationName(X, ev.Window); err != nil && name == "" {

		fmt.Printf("getApplicationName:error on window %d:\n %v\n", ev.Window, err)

		if name, err = checkQueryTreeForParent(X, ev.Window); err != nil {

			fmt.Printf("checkQueryTreeForParent:error on window %d:\n %v\n", ev.Window, err)

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
	if _, exists := allSessionOpenedAndNamedWindow[ev.Window]; !exists && name != "name-not-found" {
		allSessionOpenedAndNamedWindow[ev.Window] = name
	}
	log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
	updateWindowInfo(ev.Window, name)
}

func rootPropertyNotifyHandler(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
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
		fmt.Println("FOOOOOOOOONAAAAAAALLLLLLLLLLLYYY---ROOT")
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
		fmt.Printf("In propertyNotify:ROOT but atom was %v and wanted atom is %v or %v or %v or %v\n", ev.Atom, atom, atom1, atom2, atom3)
	}

}
