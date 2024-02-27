package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
)

func main() {

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	windows, err := getTopLevelWindows(X)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("len(windows):", len(windows))
	for _, window := range windows {
		fmt.Println(window)

		otherName, err := getApplicationName(X, window)
		if err != nil {
			fmt.Printf("\ngetApplicationName:error on window %d: %v\n\n", window, err)
		} else {
			fmt.Println(window, "====>", otherName)
		}

	}

}

func getApplicationName(X *xgbutil.XUtil, win xproto.Window) (string, error) {

	wmClass, err := xprop.PropValStrs(xprop.GetProperty(X, win, "WM_CLASS"))
	if err == nil && len(wmClass) == 2 {
		return wmClass[1], nil // The class part is typically more descriptive
	}

	// Further fallbacks or process ID resolution can go here

	return "", err // Or a more descriptive error if you prefer
}

// getTopLevelWindows returns a list of all top-level windows.
func getTopLevelWindows(X *xgbutil.XUtil) ([]xproto.Window, error) {
	return ewmh.ClientListGet(X)
}
