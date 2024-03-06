package main

import (
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

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		rootMapNotifyHandler(X, ev)
	}).Connect(X, X.RootWin())

	netActiveWindowAtom, err = xprop.Atm(X, "_NET_ACTIVE_WINDOW")
	if err != nil {
		log.Fatalf("Could not get _NET_ACTIVE_WINDOW atom: %v", err)
	}

	xevent.PropertyNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
			rootPropertyNotifyHandler(X, ev, netActiveWindowAtom)
		}).Connect(X, X.RootWin())

}

func setRootEventMask(X *xgbutil.XUtil) {

	if err = xproto.ChangeWindowAttributesChecked(X.Conn(), X.RootWin(), xproto.CwEventMask,
		[]uint32{
			xproto.EventMaskPropertyChange |
				xproto.EventMaskSubstructureNotify}).Check(); err != nil {
		log.Fatal("Failed to select notify events for root:", err)
	}

}
