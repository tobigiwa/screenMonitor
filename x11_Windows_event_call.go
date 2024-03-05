package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

func registerWindowForEvents(windowId xproto.Window) {

	xevent.DestroyNotifyFun(func(xu *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		window, ok := allCurentlyOpenedWindowMap[ev.Window]
		if !ok {
			fmt.Printf("*****name should have been in the Map, destroy*****\n\n")
			return
		}
		log.Printf("Window %d ===> %s was destroyed ev.Event:%v \n", ev.Window, window.Name, ev.Event)
		removeWindowInfo(ev.Window)
	}).Connect(X, windowId)

	xevent.FocusInFun(func(X *xgbutil.XUtil, ev xevent.FocusInEvent) {
		fmt.Printf("++++Focus in on window %d:%+v\n", ev.Event, allCurentlyOpenedWindowMap[ev.Event].Name)
		app.focusInEventHandler(ev)
	}).Connect(X, windowId)

	xevent.FocusOutFun(func(X *xgbutil.XUtil, ev xevent.FocusOutEvent) {
		fmt.Printf("-----Focus out on window %d:%+v\n", ev.Event, allCurentlyOpenedWindowMap[ev.Event].Name)
		app.focusOutEventHandler(ev)
	}).Connect(X, windowId)

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		fmt.Printf("%v\n", strings.Repeat("*", 15))
		fmt.Printf("MapNotifyFunWindow:Window:%d was mappedInWINDOW ===> allCurentlyOpenedWindowMap[ev.Window].Name:%s ---- allCurentlyOpenedWindowMap[ev.Event:%v].Name:%s---is top-level window:%v\n", ev.Window, allCurentlyOpenedWindowMap[ev.Window].Name, ev.Event, allCurentlyOpenedWindowMap[ev.Event].Name, !ev.OverrideRedirect)
		fmt.Printf("%v\n", strings.Repeat("*", 15))
	}).Connect(X, windowId)

	xevent.PropertyNotifyFun(func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		app.propertyNotifyEventHandler(X, ev)
	}).Connect(X, X.RootWin())

	// // Setup for EnterNotify and LeaveNotify events
	// xevent.EnterNotifyFun(func(X *xgbutil.XUtil, ev xevent.EnterNotifyEvent) {0
	// 	fmt.Printf("Mouse entered window %d\n", ev.Event)
	// }).Connect(X, windowId)

	// xevent.LeaveNotifyFun(func(X *xgbutil.XUtil, ev xevent.LeaveNotifyEvent) {
	// 	fmt.Printf("Mouse left window %d\n", ev.Event)
	// }).Connect(X, windowId)

	// xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
	// 	name, err := getApplicationName(X, ev.Window)
	// 	if err != nil || name == "" {
	// 		fmt.Printf("getApplicationName:error on window %d:\n %v\n", ev.Window, err)
	// 		name = "name-not-found"
	// 	}

	// 	log.Printf("Window %d ===> %s was mappedInWINDOW \n", ev.Window, name)

	// 	if tree, err := xproto.QueryTree(X.Conn(), ev.Window).Reply(); err == nil {
	// 		fmt.Printf("root of new window %v:%d is same as X.Rootwin:%v\n", name, ev.Window, tree.Root == X.RootWin())

	// 		if tree.Parent == X.RootWin() {
	// 			fmt.Println("parent is root")
	// 			// goto jump
	// 		} else {
	// 			if p, ok := allCurentlyOpenedWindowMap[tree.Parent]; ok {
	// 				fmt.Printf("parent of %v:%d is ===> %v::%v mappedInWINDOW\n", name, ev.Window, p.Name, p.ID)
	// 			} else {
	// 				fmt.Printf("parent of %v:%d is not in the mapmappedInWINDOW\n", name, ev.Window)
	// 			}
	// 			fmt.Printf("children of %v:%d are %v mappedInWINDOW\n", name, ev.Window, tree.Children)
	// 		}
	// 	} else {
	// 		fmt.Printf("error in getting tree for window %d: %v mappedInWINDOW\n", ev.Window, err)
	// 	}

	// 	// jump:
	// 	updateWindowInfo(ev.Window, name)
	// }).Connect(X, windowId)

	// xevent.UnmapNotifyFun(func(X *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
	// 	window, ok := allCurentlyOpenedWindowMap[ev.Window]
	// 	if !ok {
	// 		fmt.Printf("\n*****name should have been in the Map, unmap*****\n\n")
	// 		return
	// 	}
	// 	log.Printf("Window %d ===> %s was unmapped \n", ev.Window, window.Name)
	// }).Connect(X, windowId)

}
