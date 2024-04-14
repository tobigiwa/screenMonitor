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

var (
	netActiveWindowAtom   xproto.Atom
	netClientStackingAtom xproto.Atom
)

func registerRootWindowForEvent(X *xgbutil.XUtil) {

	xevent.DestroyNotifyFun(destroyNotifyEventFuncRoot).Connect(X, X.RootWin())

	xevent.MapNotifyFun(mapNotifyEventFuncRoot).Connect(X, X.RootWin())

	xevent.PropertyNotifyFun(func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		
		app.rootPropertyNotifyHandler(X, ev, netActiveWindowAtom, netClientStackingAtom)
	}).Connect(X, X.RootWin())
}

func setRootEventMask(X *xgbutil.XUtil) {

	err = xproto.ChangeWindowAttributesChecked(X.Conn(), X.RootWin(), xproto.CwEventMask,
		[]uint32{
			xproto.EventMaskPropertyChange |
				xproto.EventMaskStructureNotify |
				xproto.EventMaskSubstructureNotify}).Check()
	if err != nil {
		log.Fatal("Failed to select notify events for root:", err)
	}
}

func destroyNotifyEventFuncRoot(xu *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
	if window, ok := curSessionOpenedWindow[ev.Window]; ok {
		deleteWindowFromcurSessionOpenedWindowMap(ev.Window)
		log.Printf("ROOT<========Window %d:%s WAS DESTROYED!!! ev.Event:%v========>\n", ev.Window, window.Name, ev.Event)
	}
	xevent.Detach(X, ev.Window)
}

func mapNotifyEventFuncRoot(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
	fmt.Printf("\nrootMapNotifyHandler ev.window:%v ======++++++====> ev.event:%v\n", ev.Window, ev.Event)

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
	app.rootMapNotifyHandler(X, ev)
}
