package main

import (
	"LiScreMon/store"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
)

func (x11 *X11) rootMapNotifyHandler(X *xgbutil.XUtil, ev xevent.MapNotifyEvent) {
	var (
		name string
		err  error
	)

	if name, exists := curSessionNamedWindow[ev.Window]; exists {
		fmt.Printf("window:%v name resolved from window ITSELF %s:%v\n", ev.Window, name, ev.Event)
		fmt.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionOpenedWindowMap(ev.Window, name)
		return
	}

	if name, err = getWindowClassName(X, ev.Window); err != nil {
		fmt.Printf("getWindowClassName:error on window %d:\n %v\n", ev.Window, err)
		name = "name-not-found"
	}

	if name != "name-not-found" {
		fmt.Printf("Window %d ===> %s was mapped \n", ev.Window, name)
		addWindowTocurSessionOpenedWindowMap(ev.Window, name)
		addWindowTocurSessionNamedWindowMap(ev.Window, name)
	}
}

func (x11 *X11) rootPropertyNotifyHandler(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent, netActiveWindowAtom, netClientStackingAtom xproto.Atom) {

	if ev.Atom == netActiveWindowAtom {
		passedActiveWindow := netActiveWindow.WindowID

		if activeWin, err := ewmh.ActiveWindowGet(X); (err == nil) && (activeWin != 0) { // 0 is root, to much noise
			if activeWin != passedActiveWindow { // a window has become active
				s := store.ScreenTime{
					AppName: netActiveWindow.WindowName,
					Type:    store.Active,
					Time:    time.Since(netActiveWindow.TimeStamp).Hours(),
				}

				fmt.Printf("New active window ID =====> %v:%v:%v\n", activeWin, curSessionNamedWindow[activeWin], curSessionOpenedWindow[activeWin].Name)
				fmt.Printf("time elapsed for window %v:%v was %v\n", passedActiveWindow, curSessionNamedWindow[passedActiveWindow], time.Since(netActiveWindow.TimeStamp).Hours())

				if _, exists := curSessionNamedWindow[activeWin]; !exists { // if name does not already exist in curSessionNamedWindow, include it.
					if name, err := getWindowClassName(X, activeWin); err == nil {
						curSessionNamedWindow[activeWin] = name
						fmt.Printf("window:%v ====> name:%v now added in curSessionNamedWindow\n", activeWin, name)
					}
				}

				netActiveWindow.WindowID = activeWin
				netActiveWindow.WindowName = curSessionNamedWindow[activeWin]
				netActiveWindow.TimeStamp = time.Now()

				if err := x11.db.WriteUsage(s); err != nil {
					log.Fatalf("focusOutEventHandler:write to db error:%v", err)
				}

			}
		}

	}

	if ev.Atom == netClientStackingAtom {
		fmt.Printf("rootPropertyNotifyHandler:updateClientListStacking ---->---->----->on window:%v:%v\n", ev.Window, curSessionNamedWindow[ev.Window])

		windows, err := ewmh.ClientListStackingGet(X)
		if err != nil {
			log.Printf("Failed to get client list stacking: %v", err)
			return
		}

		fmt.Println("Listing stack:")
		for _, window := range windows {
			fmt.Printf("Window ID: %v, Name: %v\n", window, curSessionNamedWindow[window])
		}

		fmt.Println()

		getVisibleWindow(X.Conn(), X, windows[len(windows)-1], windows[:len(windows)-1])

	}
}

func getVisibleWindow(X *xgb.Conn, X11 *xgbutil.XUtil, activeWindow xproto.Window, otherWindows []xproto.Window) (int, error) {
	geom1, err := xproto.GetGeometry(X, xproto.Drawable(activeWindow)).Reply()
	if err != nil {
		return 0, err
	}

	rootGeom, err := xproto.GetGeometry(X, xproto.Drawable(X11.RootWin())).Reply()
	if err != nil {
		return 0, err
	}

	trans1, err := xproto.TranslateCoordinates(X, activeWindow, X11.RootWin(), geom1.X, geom1.Y).Reply()
	if err != nil {
		return 0, err
	}

	activeArea := clipToScreen(trans1, geom1, rootGeom)
	rootArea := int(rootGeom.Width) * int(rootGeom.Height)

	coverage := activeArea
	visibleWindows := []string{curSessionNamedWindow[activeWindow]}

	var prevGeom *xproto.GetGeometryReply = geom1
	var prevTrans *xproto.TranslateCoordinatesReply = trans1

	for i := len(otherWindows) - 1; i >= 0; i-- {
		otherWindow := otherWindows[i]
		geom2, err := xproto.GetGeometry(X, xproto.Drawable(otherWindow)).Reply()
		if err != nil {
			return 0, err
		}

		trans2, err := xproto.TranslateCoordinates(X, otherWindow, X11.RootWin(), geom2.X, geom2.Y).Reply()
		if err != nil {
			return 0, err
		}

		area2 := clipToScreen(trans2, geom2, rootGeom)

		x_overlap := max(0, min(int(prevTrans.DstX)+int(prevGeom.Width), int(trans2.DstX)+int(geom2.Width))-max(int(prevTrans.DstX), int(trans2.DstX)))
		y_overlap := max(0, min(int(prevTrans.DstY)+int(prevGeom.Height), int(trans2.DstY)+int(geom2.Height))-max(int(prevTrans.DstY), int(trans2.DstY)))

		overlapArea := x_overlap * y_overlap

		overlapRatio := float64(overlapArea) / float64(area2)

		fmt.Println(float64(overlapArea), float64(area2))

		if overlapRatio < 0.97 {
			coverage += area2 - overlapArea
			visibleWindows = append(visibleWindows, curSessionNamedWindow[otherWindow])
		}

		fmt.Printf("Windows --> %v, last Overlap: %v, stage: %v%%\n", visibleWindows, overlapRatio, int(math.Round(float64(coverage)/float64(rootArea)*100)))

		if float64(coverage)/float64(rootArea) >= 0.80 {
			fmt.Printf("The combination of windows %v covers more than 80%% of the screen\n", visibleWindows)
			return 100, nil
		}

		prevGeom = geom2
		prevTrans = trans2
	}

	return 0, nil // None of the windows in the slice caused the total coverage to be more than 80% of the screen
}

func clipToScreen(trans *xproto.TranslateCoordinatesReply, geom *xproto.GetGeometryReply, rootGeom *xproto.GetGeometryReply) int {
	x := int(trans.DstX)
	y := int(trans.DstY)
	width := int(geom.Width)
	height := int(geom.Height)
	screenWidth := int(rootGeom.Width)
	screenHeight := int(rootGeom.Height)

	if x < 0 {
		width += x
		x = 0
	}
	if y < 0 {
		height += y
		y = 0
	}
	if x+width > screenWidth {
		width = screenWidth - x
	}
	if y+height > screenHeight {
		height = screenHeight - y
	}
	return max(0, width) * max(0, height) // Return the area
}

// func getVisibleWindow(X *xgb.Conn, X11 *xgbutil.XUtil, activeWindow xproto.Window, otherWindows []xproto.Window) (int, error) {
// 	geom1, err := xproto.GetGeometry(X, xproto.Drawable(activeWindow)).Reply()
// 	if err != nil {
// 		return 0, err
// 	}

// 	rootGeom, err := xproto.GetGeometry(X, xproto.Drawable(X11.RootWin())).Reply()
// 	if err != nil {
// 		return 0, err
// 	}

// 	trans1, err := xproto.TranslateCoordinates(X, activeWindow, X11.RootWin(), geom1.X, geom1.Y).Reply()
// 	if err != nil {
// 		return 0, err
// 	}

// 	// activeArea := int(geom1.Width) * int(geom1.Height)
// 	activeArea := clipToScreen(trans1, geom1, rootGeom)
// 	rootArea := int(rootGeom.Width) * int(rootGeom.Height)

// 	coverage := activeArea
// 	visibleWindows := []string{curSessionNamedWindow[activeWindow]}

// 	for i := len(otherWindows) - 1; i >= 0; i-- {
// 		otherWindow := otherWindows[i]
// 		geom2, err := xproto.GetGeometry(X, xproto.Drawable(otherWindow)).Reply()
// 		if err != nil {
// 			return 0, err
// 		}

// 		trans2, err := xproto.TranslateCoordinates(X, otherWindow, X11.RootWin(), geom2.X, geom2.Y).Reply()
// 		if err != nil {
// 			return 0, err
// 		}

// 		// area2 := int(geom2.Width) * int(geom2.Height)
// 		area2 := clipToScreen(trans2, geom2, rootGeom)

// 		x_overlap := max(0, min(int(trans1.DstX)+int(geom1.Width), int(trans2.DstX)+int(geom2.Width))-max(int(trans1.DstX), int(trans2.DstX)))
// 		y_overlap := max(0, min(int(trans1.DstY)+int(geom1.Height), int(trans2.DstY)+int(geom2.Height))-max(int(trans1.DstY), int(trans2.DstY)))

// 		overlapArea := x_overlap * y_overlap

// 		overlapRatio := float64(overlapArea) / float64(area2)
// 		fmt.Println(float64(overlapArea), float64(area2))

// 		if overlapRatio < 0.97 {
// 			coverage += area2 - overlapArea
// 			visibleWindows = append(visibleWindows, curSessionNamedWindow[otherWindow])
// 		}

// 		fmt.Printf("Windows --> %v, last Overlap: %v, stage: %v%%\n", visibleWindows, overlapRatio, int(math.Round(float64(coverage)/float64(rootArea)*100)))

// 		if float64(coverage)/float64(rootArea) >= 0.80 {
// 			fmt.Printf("The combination of windows %v covers more than 80%% of the screen\n", visibleWindows)
// 			return 100, nil
// 		}
// 	}

// 	return 0, nil // None of the windows in the slice caused the total coverage to be more than 80% of the screen
// }
