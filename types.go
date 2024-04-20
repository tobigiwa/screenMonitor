package main

import (
	"LiScreMon/store"
	"sync"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

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
