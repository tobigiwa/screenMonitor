package main

import (
	"LiScreMon/store"
	"sync"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

var (
	allCurentlyOpenedWindowMap     = make(map[xproto.Window]WindowInfo)
	allSessionOpenedAndNamedWindow = make(map[xproto.Window]string)
)

type WindowInfo struct {
	ID   xproto.Window
	Name string
}
type DoNotCopy [0]sync.Mutex
type focusEvent struct {
	WindowID   xproto.Window
	AppName    string
	ScreenType string
	Time       time.Time
	DoNotCopy
}

type X11 struct {
	db store.Repository
}
