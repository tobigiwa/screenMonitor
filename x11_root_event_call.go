package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
)

var (
	netActiveWindowAtom xproto.Atom
)

func registerRootWindowForEvent(X *xgbutil.XUtil) {

	xevent.DestroyNotifyFun(func(xu *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		if window, ok := curSessionOpenedWindow[ev.Window]; ok {
			deleteWindowFromcurSessionOpenedWindowMap(ev.Window)
			log.Printf("ROOT<========Window %d:%s was destroyed ev.Event:%v========>\n", ev.Window, window.Name, ev.Event)
		}
		xevent.Detach(X, ev.Window)
	}).Connect(X, X.RootWin())

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		fmt.Printf("\nrootMapNotifyHandler ev.window:%v ======++++++====> ev.event:%v\n", ev.Window, ev.Event)
		app.rootMapNotifyHandler(X, ev)
	}).Connect(X, X.RootWin())

	xevent.VisibilityNotifyFun(func(xu *xgbutil.XUtil, ev xevent.VisibilityNotifyEvent) {
		fmt.Printf("\nrootVisibilityNotifyHandler:::window:%v:name-%v got visibility of state--- %v\n", ev.Window, curSessionNamedWindow[ev.Window], ev.State)
	}).Connect(X, X.RootWin())

	netActiveWindowAtom, err = xprop.Atm(X, "_NET_ACTIVE_WINDOW")
	if err != nil {
		log.Fatalf("Could not get _NET_ACTIVE_WINDOW atom: %v", err)
	}

	xevent.PropertyNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
			app.rootPropertyNotifyHandler(X, ev, netActiveWindowAtom)
		}).Connect(X, X.RootWin())
}

func setRootEventMask(X *xgbutil.XUtil) {

	err = xproto.ChangeWindowAttributesChecked(X.Conn(), X.RootWin(), xproto.CwEventMask,
		[]uint32{
			xproto.EventMaskPropertyChange |
				xproto.EventMaskVisibilityChange |
				xproto.EventMaskStructureNotify |
				xproto.EventMaskSubstructureNotify}).Check()
	if err != nil {
		log.Fatal("Failed to select notify events for root:", err)
	}

}
