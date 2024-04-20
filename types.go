package main

import (
	"LiScreMon/store"
	"sync"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

var (
	// curSessionNamedWindow is a map of all current session "named" windows.
	// An X session is typically a time between login and logout (or restart/shutdown).
	// Only windows with knowm WM_CLASS are added to this map. The X_ID are always unique
	// for a particular window in each session.
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

type DoNotCopy [0]sync.Mutex

type X11 struct {
	db store.IRepository
}
