package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
)

func SetDefaultFocusWindow(X *xgbutil.XUtil) error {
	FocusWindow, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		return err
	}
	globalFocusEvent.WindowID = FocusWindow
	return nil
}

// Function to remove a window from the map (e.g., when it's closed)
func removeWindowInfo(win xproto.Window) {
	delete(allCurentlyOpenedWindowMap, win)
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

// Function to update the window map with a new or existing window
func updateWindowInfo(windowID xproto.Window, name string) {
	if _, exists := allCurentlyOpenedWindowMap[windowID]; !exists {
		allCurentlyOpenedWindowMap[windowID] = WindowInfo{
			ID:   windowID,
			Name: name,
		}

		xproto.ChangeWindowAttributes(X.Conn(), windowID, xproto.CwEventMask,
			[]uint32{
				xproto.EventMaskFocusChange |
					// xproto.EventMaskPropertyChange |
					// xproto.EventMaskEnterWindow |
					// xproto.EventMaskLeaveWindow |
					xproto.EventMaskSubstructureNotify,
			})
		registerWindowForEvents(windowID)
		return
	}
	fmt.Printf("\nwindow already present, weird\n")
}

func getApplicationName(X *xgbutil.XUtil, win xproto.Window) (string, error) {

	var (
		wmName  string
		err1    error
		err2    error
		err3    error
		err4    error
		err5    error
		err6    error
		err7    error
		wmClass []string
		wmPID   uint
		cmdline []byte
	)

	if wmClass, err1 = xprop.PropValStrs(xprop.GetProperty(X, win, "WM_CLASS")); err1 == nil && (len(wmClass) == 2) {
		fmt.Printf("\nthe wmClass array is on windowID:%v -----> %+v\n", win, wmClass)
		return wmClass[1], nil // The class part is typically more descriptive
	}

	if wmName, err2 = ewmh.WmNameGet(X, win); err2 == nil && wmName != "" {
		fmt.Printf("the WmNameGet on windowID:%v -----> %+v\n", win, wmName)
		return wmName, nil
	}

	if wmName, err7 = xprop.PropValStr(xprop.GetProperty(X, win, "WM_NAME")); err1 == nil && wmName != "" {
		fmt.Printf("\nthe wmName from deprecated is on windowID:%v -----> %+v\n", win, wmClass)
		return wmName, nil // The class part is typically more descriptive
	}

	if wmName, err3 = ewmh.WmVisibleNameGet(X, win); err2 == nil && wmName != "" {
		fmt.Printf("the WmVisibleNameGet on windowID:%v -----> %+v\n", win, wmName)
		return wmName, nil
	}

	if wmName, err4 = ewmh.WmIconNameGet(X, win); err2 == nil && wmName != "" {
		fmt.Printf("the WmIconNameGet on windowID:%v -----> %+v\n", win, wmName)
		return wmName, nil
	}

	if wmName, err5 = ewmh.WmVisibleIconNameGet(X, win); err2 == nil && wmName != "" {
		fmt.Printf("the WmVisibleIconNameGet on windowID:%v -----> %+v\n", win, wmName)
		return wmName, nil
	}

	if wmPID, err6 = ewmh.WmPidGet(X, win); err3 == nil {
		fmt.Printf("the wmPID on windowID:%v -----> %+v\n", win, wmPID)
		if cmdline, err3 = os.ReadFile(fmt.Sprintf("/proc/%d/comm", wmPID)); err3 == nil {
			return strings.ReplaceAll(string(cmdline), "\x00", " "), nil
		}
	}

	// Further fallbacks or process ID resolution can go here

	return "", fmt.Errorf("cannot resolve window name %v:%v:%v:%v:%v:%v:%v", err1, err2, err3, err4, err5, err6, err7)
}

func checkQueryTreeForParent(X *xgbutil.XUtil, window xproto.Window) (string, error) {

	var (
		tree *xproto.QueryTreeReply
		err  error
	)

	if tree, err = xproto.QueryTree(X.Conn(), window).Reply(); err != nil {
		return "", err
	}

	if parentName, ok := allSessionOpenedAndNamedWindow[tree.Parent]; ok {
		fmt.Printf("window:%v name resolved from parent %s:\n", window, parentName)
		return parentName, nil
	}

	if windowName, ok := allSessionOpenedAndNamedWindow[window]; ok {
		fmt.Printf("window:%v name resolved from window itself %s:\n", window, windowName)
		return windowName, nil
	}

	return "", fmt.Errorf("window %d not found in allSessionOpenedAndNamedWindow:window is top-level== %v", window, tree.Parent == X.RootWin())
}
