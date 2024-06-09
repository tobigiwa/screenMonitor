package monitoring

import (
	"fmt"
	"log"
	"pkg/types"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
)

func rootMapNotifyHandler(x11Conn *xgbutil.XUtil, ev xevent.MapNotifyEvent) {

	name, _ := getWindowClassName(x11Conn, ev.Window)
	registerWindowForEvents(ev.Window) // For DestroyNotify on the window

	if name != "" {
		addWindowTocurSessionNamedWindowMap(ev.Window, name)

	}
}

func (x11 *X11Monitor) rootPropertyNotifyHandler(x11Conn *xgbutil.XUtil, ev xevent.PropertyNotifyEvent, netActiveWindowAtom, netClientStackingAtom xproto.Atom) {

	if ev.Atom == netActiveWindowAtom {
		if currActiveWindow, err := ewmh.ActiveWindowGet(x11Conn); (err == nil) && (currActiveWindow != 0) { // 0 is root, to much noise

			if formerActiveWindow := netActiveWindow; formerActiveWindow.WindowID != currActiveWindow { // this helps takes care of noise from tabs switch

				x11.windowChangeCh <- struct{}{}

				if formerActiveWindow.WindowID == xevent.NoWindow { // at first run i.e on boot
					netActiveWindow.WindowID = currActiveWindow                                   // SET THE WINDOW ID
					netActiveWindow.TimeStamp = time.Now()                                        // SET THE TIME
					netActiveWindow.WindowName, _ = getWindowClassName(x11Conn, currActiveWindow) // SET THE NAME
					return
				}

				if formerActiveWindow.WindowName == "" { // this might be a not so needed check, cos -->
					formerActiveWindow.WindowName, _ = getWindowClassName(x11Conn, currActiveWindow) // NET_ACTIVE_WINDOW SHOULD ALWAYS HAVE A NAME, if not, that is lost metric
				}

				s := types.ScreenTime{
					WindowID: formerActiveWindow.WindowID,
					AppName:  formerActiveWindow.WindowName,
					Type:     types.Active,
					Duration: time.Since(formerActiveWindow.TimeStamp).Hours(),
					Interval: types.TimeInterval{Start: formerActiveWindow.TimeStamp, End: time.Now()},
				}

				fmt.Printf("New active window ID =====> %v:%v\ntime elapsed for last window %v:%v was %vsecs\n",
					currActiveWindow, curSessionNamedWindow[currActiveWindow], formerActiveWindow.WindowID, curSessionNamedWindow[formerActiveWindow.WindowID], time.Since(netActiveWindow.TimeStamp).Seconds())

				// SETTING THE NEW _NET_ACTIVE_WINDOW
				netActiveWindow.WindowID = currActiveWindow                                   // SET THE WINDOW ID
				netActiveWindow.TimeStamp = time.Now()                                        // SET THE TIME
				netActiveWindow.WindowName, _ = getWindowClassName(x11Conn, currActiveWindow) // SET THE NAME

				if s.AppName != "" { // AS mentioned earlier, if we don't have a name, lost metric
					if err := x11.Db.WriteUsage(s); err != nil {
						log.Fatalf("write to db error:%v", err)
					}
				}
			}
		}
	}
}
