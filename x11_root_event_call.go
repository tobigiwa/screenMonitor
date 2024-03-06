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

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		fmt.Printf("\nrootMapNotifyHandler ev.window:%v ======++++++====> ev.event:%v\n", ev.Window, ev.Event)
		rootMapNotifyHandler(X, ev)
	}).Connect(X, X.RootWin())

	netActiveWindowAtom, err = xprop.Atm(X, "_NET_ACTIVE_WINDOW")
	if err != nil {
		log.Fatalf("Could not get _NET_ACTIVE_WINDOW atom: %v", err)
	}

	// attr, err := xproto.GetWindowAttributes(X.Conn(), X.RootWin()).Reply()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // // Remove the MotionNotify mask
	// _ = attr.AllEventMasks & ^(uint32)(xproto.EventMaskPropertyChange)

	xevent.PropertyNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
			fmt.Printf("\nrootPropertyNotifyHandler ev.window:%v ======++++++====> got atom %v, expecting atom %v\n\n", ev.Window, ev.Atom, netActiveWindowAtom)
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
