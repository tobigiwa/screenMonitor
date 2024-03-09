package main

import (
	"LiScreMon/store"
	"sync"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

var (
	// curSessionOpenedWindow is a map of all currently opened windows.
	// if a window is closed, it should be removed from this map.
	// An X session is typically a time between login and logout (or restart/shutdown).
	// This map is update and deleted with the
	/*
		func addWindowTocurSessionOpenedWindowMap
		func deleteWindowFromcurSessionOpenedWindowMap
	*/
	curSessionOpenedWindow = make(map[xproto.Window]WindowInfo)

	// curSessionNamedWindow is a map of all current session "named" windows.
	// An X session is typically a time between login and logout (or restart/shutdown).
	// Only windows with knowm WM_CLASS are added to this map.
	curSessionNamedWindow = make(map[xproto.Window]string)
)

type WindowInfo struct {
	ID   xproto.Window
	Name string
}

type netActiveWindowInfo struct {
	WindowID   xproto.Window
	WindowName string
	TimeStamp  time.Time
	DoNotCopy
}

var netActiveWindow = &netActiveWindowInfo{}

type DoNotCopy [0]sync.Mutex

type X11 struct {
	db store.Repository
}
