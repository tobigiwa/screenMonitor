package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
)

func getInitActiveWindow(X *xgbutil.XUtil) error {
	FocusWindow, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		return err
	}
	globalFocusEvent.WindowID = FocusWindow
	globalFocusEvent.Time = time.Now()

	return nil
}

// func InitMonitoringEvent(X *xgbutil.XUtil, windowIDs []xproto.Window) {
// 	for _, windowId := range windowIDs {
// 		registerWindowForEvents(windowId)
// 	}
// }

// currentlyOpenedWindows returns a list of all top-level windows.
func currentlyOpenedWindows(X *xgbutil.XUtil) ([]xproto.Window, error) {
	return ewmh.ClientListGet(X)
}

// deleteWindowInfo adds to the
/* curSessionOpenedWindow map */
// and also set Event mask on newly added windows
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

		xproto.ChangeWindowAttributes(X.Conn(), windowID, xproto.CwEventMask,
			[]uint32{
				xproto.EventMaskFocusChange |
					// xproto.EventMaskPropertyChange |
					xproto.EventMaskSubstructureNotify,
			})

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

func getApplicationName(X *xgbutil.XUtil, win xproto.Window) (string, error) {

	var (
		err error
	)

	if wmClass, err := xprop.PropValStrs(xprop.GetProperty(X, win, "WM_CLASS")); err == nil && (len(wmClass) == 2) {
		return wmClass[1], nil
	}

	return "", fmt.Errorf("cannot resolve window name %v", err)
}

func checkQueryTreeForParent(X *xgbutil.XUtil, window xproto.Window) (string, error) {

	var (
		tree *xproto.QueryTreeReply
		err  error
	)

	if tree, err = xproto.QueryTree(X.Conn(), window).Reply(); err != nil {
		return "", err
	}

	if parentName, ok := curSessionNamedWindow[tree.Parent]; ok {
		fmt.Printf("window:%v name resolved from parent %s:\n", window, parentName)
		return parentName, nil
	}

	if windowName, ok := curSessionNamedWindow[window]; ok {
		fmt.Printf("window:%v name resolved from window itself %s:\n", window, windowName)
		return windowName, nil
	}

	return "", fmt.Errorf("window %v not found in curSessionNamedWindow, window is top-level == %v", window, tree.Parent == X.RootWin())
}
