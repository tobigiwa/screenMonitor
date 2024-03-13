package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
)

func InitNetActiveWindow(X *xgbutil.XUtil) error {
	activeWin, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		return err
	}

	netActiveWindow.WindowID = activeWin
	netActiveWindow.WindowName = curSessionNamedWindow[activeWin]
	netActiveWindow.TimeStamp = time.Now()

	return nil
}

func InitMonitoringEvent(X *xgbutil.XUtil, windowIDs []xproto.Window) {
	for _, windowId := range windowIDs {
		registerWindowForEvents(windowId)
	}
}

// currentlyOpenedWindows returns a list of all top-level windows.
func currentlyOpenedWindows(X *xgbutil.XUtil) ([]xproto.Window, error) {
	return ewmh.ClientListGet(X)
}

// deleteWindowInfo deletes from the
/* curSessionOpenedWindow map */
func deleteWindowFromcurSessionOpenedWindowMap(win xproto.Window) {
	delete(curSessionOpenedWindow, win)
}

// addWindowTocurSessionOpenedWindowMap adds to the
/* curSessionOpenedWindow map */
// set Event mask on newly added windows and register them for events.
func addWindowTocurSessionOpenedWindowMap(windowID xproto.Window, name string) {
	if _, exists := curSessionOpenedWindow[windowID]; !exists {
		curSessionOpenedWindow[windowID] = WindowInfo{
			ID:   windowID,
			Name: name,
		}

		err := xproto.ChangeWindowAttributesChecked(X.Conn(), windowID, xproto.CwEventMask,
			[]uint32{
				xproto.EventMaskStructureNotify |
					xproto.EventMaskSubstructureNotify}).Check()
		if err != nil {
			log.Fatalf("Failed to select notify events for window:%v:%v: error: %v", windowID, name, err)
		}

		registerWindowForEvents(windowID)
		return
	}
}

// addWindowTocurSessionNamedWindowMap adds to the
/* curSessionNamedWindow map */
// there is no need to set Event mask on newly added windows and widows are not deleted from this map.
// This map is used to resolve the name of the known windows.
func addWindowTocurSessionNamedWindowMap(windowID xproto.Window, name string) {
	if _, exists := curSessionNamedWindow[windowID]; !exists {
		curSessionNamedWindow[windowID] = name
		return
	}
}

func getWindowClassName(X *xgbutil.XUtil, win xproto.Window) (string, error) {

	wmClass, err1 := xprop.PropValStrs(xprop.GetProperty(X, win, "WM_CLASS"))
	if err1 == nil && (len(wmClass) == 2) {
		return wmClass[1], nil
	}

	name, err2 := checkQueryTreeForParent(X, win)
	if err2 == nil && (name != "") {
		return name, nil
	}

	return "", fmt.Errorf("error on resolving name for window %d: %v, %w", win, err1, err2)
}

func checkQueryTreeForParent(X *xgbutil.XUtil, window xproto.Window) (string, error) {

	var (
		tree *xproto.QueryTreeReply
		err  error
	)

	if tree, err = xproto.QueryTree(X.Conn(), window).Reply(); err == nil {
		if parentName, ok := curSessionNamedWindow[tree.Parent]; ok && (tree.Parent != tree.Root) { // root would'nt have been added to curSessionNamedWindow
			return parentName, nil
		}

		ChildrenLen := len(tree.Children)
		for i := 0; i < ChildrenLen; i++ {
			if childName, ok := curSessionNamedWindow[tree.Children[i]]; ok { // noticed this behavior from vscode
				return childName, nil
			}

		}
	}
	return "", err
}

func WmTransientForGet(xu *xgbutil.XUtil, win xproto.Window) (xproto.Window, error) {
	raw, err := xprop.GetProperty(xu, win, "WM_TRANSIENT_FOR")
	if err != nil {
		return 0, err
	}
	return xprop.PropValWindow(raw, err)
}


