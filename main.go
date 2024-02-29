package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
)

var (
	X         *xgbutil.XUtil
	err       error
	windowMap = make(map[xproto.Window]WindowInfo)
)

type WindowInfo struct {
	ID   xproto.Window
	Name string
	TimeActivity WindowTimeActivity
}
type WindowTimeActivity struct {
    Duration time.Time
    TimeStore time.Time
}

func main() {

	if X, err = xgbutil.NewConn(); err != nil {
		log.Fatal(err)
	}

	if err = xproto.ChangeWindowAttributesChecked(X.Conn(), X.RootWin(), xproto.CwEventMask,
		[]uint32{
			xproto.EventMaskFocusChange |
				xproto.EventMaskLeaveWindow |
				xproto.EventMaskEnterWindow |
				xproto.EventMaskSubstructureNotify}).Check(); err != nil {
		log.Fatal("Failed to select notify events:", err)
	}

	// xproto.QueryTree()
	// xwindow.New(X, X.RootWin()).Listen(xproto.EventMaskSubstructureNotify)

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		name, err := getApplicationName(X, ev.Window)
		if err != nil || name == "" {
			fmt.Printf("getApplicationName:error on window %d:\n %v\n", ev.Window, err)
			name = "name-not-found"
		}

		log.Printf("Window %d ===> %s was mapped \n", ev.Window, name)

		if tree, err := xproto.QueryTree(X.Conn(), ev.Window).Reply(); err == nil {
			fmt.Printf("root of new window %v:%d is same as X.Rootwin:%v\n", name, ev.Window, tree.Root == X.RootWin())

			if tree.Parent == X.RootWin() {
				fmt.Println("parent is root")
				// goto jump
			} else {
				if p, ok := windowMap[tree.Parent]; ok {
					fmt.Printf("parent of %v:%d is ===> %v::%v\n", name, ev.Window, p.Name, p.ID)
				} else {
					fmt.Printf("parent of %v:%d is not in the map\n", name, ev.Window)
				}
				fmt.Printf("children of %v:%d are %v\n", name, ev.Window, tree.Children)
			}
		} else {
			fmt.Printf("error in getting tree for window %d: %v\n", ev.Window, err)
		}

		// jump:
		updateWindowInfo(ev.Window, name)
	}).Connect(X, X.RootWin())

	windows, err := getTopLevelWindows(X)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("len(windows):", len(windows))
	for _, window := range windows {
		otherName, err := getApplicationName(X, window)
		if err != nil {
			fmt.Printf("\ngetApplicationName:error on window %d: %v\n\n", window, err)
		} else {
			fmt.Println(window, "====>", otherName)
		}

		updateWindowInfo(window, otherName)
	}

	InitMonitoringEvent(X, windows)

	xevent.Main(X)
}

func getApplicationName(X *xgbutil.XUtil, win xproto.Window) (string, error) {

	var (
		wmName  string
		err1    error
		err2    error
		err3    error
		wmClass []string
		wmPID   uint
		cmdline []byte
	)

	if wmClass, err1 = xprop.PropValStrs(xprop.GetProperty(X, win, "WM_CLASS")); err1 == nil && (len(wmClass) == 2) {
		fmt.Printf("the wmClass array is on windowID:%v -----> %+v\n", win, wmClass)
		return wmClass[1], nil // The class part is typically more descriptive
	}

	if wmName, err2 = ewmh.WmNameGet(X, win); err2 == nil && wmName != "" {
		fmt.Printf("the wmName on windowID:%v -----> %+v\n", win, wmName)
		return wmName, nil
	}

	if wmPID, err3 = ewmh.WmPidGet(X, win); err3 == nil {
		fmt.Printf("the wmPID on windowID:%v -----> %+v\n", win, wmPID)
		if cmdline, err3 = os.ReadFile(fmt.Sprintf("/proc/%d/comm", wmPID)); err3 == nil {
			return strings.ReplaceAll(string(cmdline), "\x00", " "), nil
		}
	}

	// Further fallbacks or process ID resolution can go here

	return "", fmt.Errorf("cannot reolve window name %w:%w:%w", err1, err2, err3)
}

// getTopLevelWindows returns a list of all top-level windows.
func getTopLevelWindows(X *xgbutil.XUtil) ([]xproto.Window, error) {
	return ewmh.ClientListGet(X)
}

// Function to update the window map with a new or existing window
func updateWindowInfo(windowID xproto.Window, name string) {
	if _, exists := windowMap[windowID]; !exists {
		windowMap[windowID] = WindowInfo{
			ID:   windowID,
			Name: name,
		}

		xproto.ChangeWindowAttributes(X.Conn(), windowID, xproto.CwEventMask,
			[]uint32{
				xproto.EventMaskFocusChange |
					xproto.EventMaskEnterWindow |
					xproto.EventMaskLeaveWindow |
					xproto.EventMaskSubstructureNotify,
			})
		monitorWindow(windowID)
		return
	}
	fmt.Printf("\nwindow already present, weird\n")
}

// Function to remove a window from the map (e.g., when it's closed)
func removeWindowInfo(win xproto.Window) {
	delete(windowMap, win)
}

func InitMonitoringEvent(X *xgbutil.XUtil, windowIDs []xproto.Window) {
	for _, windowId := range windowIDs {
		monitorWindow(windowId)
	}
}

func monitorWindow(windowId xproto.Window) {

	xevent.DestroyNotifyFun(func(xu *xgbutil.XUtil, ev xevent.DestroyNotifyEvent) {
		window, ok := windowMap[ev.Window]
		if !ok {
			fmt.Printf("*****name should have been in the Map, destroy*****\n\n")
			return
		}
		log.Printf("Window %d ===> %s was destroyed ev.Event:%v \n", ev.Window, window.Name, ev.Event)
		removeWindowInfo(ev.Window)
	}).Connect(X, windowId)

	xevent.FocusInFun(func(X *xgbutil.XUtil, ev xevent.FocusInEvent) {
		fmt.Printf("Focus in on window %d\n", ev.Event)
	}).Connect(X, windowId)

	xevent.FocusOutFun(func(X *xgbutil.XUtil, ev xevent.FocusOutEvent) {
		fmt.Printf("Focus out on window %d\n", ev.Event)
	}).Connect(X, windowId)

	// Setup for EnterNotify and LeaveNotify events
	xevent.EnterNotifyFun(func(X *xgbutil.XUtil, ev xevent.EnterNotifyEvent) {
		fmt.Printf("Mouse entered window %d\n", ev.Event)
	}).Connect(X, windowId)

	xevent.LeaveNotifyFun(func(X *xgbutil.XUtil, ev xevent.LeaveNotifyEvent) {
		fmt.Printf("Mouse left window %d\n", ev.Event)
	}).Connect(X, windowId)

	xevent.MapNotifyFun(func(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
		name, err := getApplicationName(X, ev.Window)
		if err != nil || name == "" {
			fmt.Printf("getApplicationName:error on window %d:\n %v\n", ev.Window, err)
			name = "name-not-found"
		}

		log.Printf("Window %d ===> %s was mappedInWINDOW \n", ev.Window, name)

		if tree, err := xproto.QueryTree(X.Conn(), ev.Window).Reply(); err == nil {
			fmt.Printf("root of new window %v:%d is same as X.Rootwin:%v\n", name, ev.Window, tree.Root == X.RootWin())

			if tree.Parent == X.RootWin() {
				fmt.Println("parent is root")
				// goto jump
			} else {
				if p, ok := windowMap[tree.Parent]; ok {
					fmt.Printf("parent of %v:%d is ===> %v::%v mappedInWINDOW\n", name, ev.Window, p.Name, p.ID)
				} else {
					fmt.Printf("parent of %v:%d is not in the mapmappedInWINDOW\n", name, ev.Window)
				}
				fmt.Printf("children of %v:%d are %v mappedInWINDOW\n", name, ev.Window, tree.Children)
			}
		} else {
			fmt.Printf("error in getting tree for window %d: %v mappedInWINDOW\n", ev.Window, err)
		}

		// jump:
		updateWindowInfo(ev.Window, name)
	}).Connect(X, windowId)

	// xevent.UnmapNotifyFun(func(X *xgbutil.XUtil, ev xevent.UnmapNotifyEvent) {
	// 	window, ok := windowMap[ev.Window]
	// 	if !ok {
	// 		fmt.Printf("\n*****name should have been in the Map, unmap*****\n\n")
	// 		return
	// 	}
	// 	log.Printf("Window %d ===> %s was unmapped \n", ev.Window, window.Name)
	// }).Connect(X, windowId)

}
