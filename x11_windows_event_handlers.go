package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/xgbutil/xevent"
)

var globalFocusEvent focusEvent

func (x11 X11) focusInEventHandler(fv xevent.FocusInEvent) {

	if fv.Event != globalFocusEvent.WindowID {
		globalFocusEvent.Time = time.Now()
		globalFocusEvent.WindowID = fv.Event

	}
}

func (x11 X11) focusOutEventHandler(fv xevent.FocusOutEvent) {

	var duration time.Duration
	currentFocusWIndow := globalFocusEvent.WindowID
	if fv.Event == currentFocusWIndow {
		duration = time.Since(globalFocusEvent.Time)

		fmt.Printf("    time elapsed was %v\n", duration.Minutes())

		// s := store.ScreenTime{
		// 	AppName: globalFocusEvent.AppName,
		// 	Type:    store.Active,
		// 	Time:    duration.Minutes(),
		// }

		// if err := x11.db.WriteUsuage(s); err != nil {
		// 	log.Fatalf("focusOutEventHandler:write to db error:%v", err)
		// }
	}
}
