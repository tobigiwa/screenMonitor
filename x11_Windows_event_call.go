package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
)

func registerWindowForEvents(windowId xproto.Window) {

	xevent.DestroyNotifyFun(func(xu *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		if window, ok := curSessionOpenedWindow[ev.Window]; ok {
			deleteWindowFromcurSessionOpenedWindowMap(ev.Window)
			log.Printf("WINDOW<========Window %d:%s WAS DESTROYED!!! ev.Event:%v========>\n", ev.Window, window.Name, ev.Event)
		}
		xevent.Detach(X, windowId)
	}).Connect(X, windowId)

	

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		fmt.Printf("\nWINDOWMapNotifyHandler ev.window:%v ======++++++====> ev.event:%v\n", ev.Window, ev.Event)
		if ev.OverrideRedirect {
			return
		}

		if transientFor, err := xprop.PropValWindow(xprop.GetProperty(X, ev.Window, "WM_TRANSIENT_FOR")); err == nil && transientFor != 0 {
			fmt.Println("This window is transient for window", transientFor)
			return
		}

		if windowTypes, err := ewmh.WmWindowTypeGet(X, ev.Window); err == nil {
			for i := 0; i < len(windowTypes); i++ {
				if windowTypes[i] == "_NET_WM_WINDOW_TYPE_SPLASH" {
					fmt.Println("This is a splash window")
					fmt.Printf("the array was %+v\n", windowTypes)
					return
				}
			}
		}

		app.windowMapNotifyHandler(X, ev)
	}).Connect(X, windowId)

	xevent.VisibilityNotifyFun(func(xu *xgbutil.XUtil, ev xevent.VisibilityNotifyEvent) {
		switch ev.State {
		case xproto.VisibilityUnobscured:
			fmt.Printf("\nWINDOW::::window:%v:name-%v got visibility of state--- %v: VisibilityUnobscured\n", ev.Window, curSessionNamedWindow[ev.Window], ev.State)
		case xproto.VisibilityPartiallyObscured:
			fmt.Printf("\nWINDOW::::window:%v:name-%v got visibility of state--- %v:VisibilityPartiallyObscured\n", ev.Window, curSessionNamedWindow[ev.Window], ev.State)
		case xproto.VisibilityFullyObscured:
			fmt.Printf("\nWINDOW::::window:%v:name-%v got visibility of state--- %v: VisibilityFullyObscured\n", ev.Window, curSessionNamedWindow[ev.Window], ev.State)
		}
	}).Connect(X, windowId)

}
