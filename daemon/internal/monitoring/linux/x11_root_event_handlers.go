package monitoring

import (
	"fmt"
	"log"
	"time"

	"LiScreMon/daemon/internal/database/repository"

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

				if formerActiveWindow.WindowID == xevent.NoWindow { // at first run
					netActiveWindow.WindowID = currActiveWindow                                   // SET THE WINDOW ID
					netActiveWindow.TimeStamp = time.Now()                                        // SET THE TIME
					netActiveWindow.WindowName, _ = getWindowClassName(x11Conn, currActiveWindow) // SET THE NAME

					if netActiveWindow.WindowName != "" {
						curSessionNamedWindow[currActiveWindow] = netActiveWindow.WindowName // include it the named windows
					}
					return
				}

				if formerActiveWindow.WindowName == "" {
					formerActiveWindow.WindowName, _ = getWindowClassName(x11Conn, currActiveWindow) // NET_ACTIVE_WINDOW SHOULD ALWAYS HAVE A NAME, if not, that is lost metric.
				}

				s := repository.ScreenTime{
					WindowID: formerActiveWindow.WindowID,
					AppName:  formerActiveWindow.WindowName,
					Type:     repository.Active,
					Duration: time.Since(formerActiveWindow.TimeStamp).Hours(),
					Interval: repository.TimeInterval{Start: formerActiveWindow.TimeStamp, End: time.Now()},
				}

				fmt.Printf("New active window ID =====> %v:%v\ntime elapsed for last window %v:%v was %vsecs\n",
					currActiveWindow, curSessionNamedWindow[currActiveWindow], formerActiveWindow.WindowID, curSessionNamedWindow[formerActiveWindow.WindowID], time.Since(netActiveWindow.TimeStamp).Seconds())

				// SETTING THE NEW _NET_ACTIVE_WINDOW
				var ok bool
				netActiveWindow.WindowID = currActiveWindow                                        // SET THE WINDOW ID
				netActiveWindow.TimeStamp = time.Now()                                             // SET THE TIME
				if netActiveWindow.WindowName, ok = curSessionNamedWindow[currActiveWindow]; !ok { // SET THE NAME
					netActiveWindow.WindowName, _ = getWindowClassName(x11Conn, currActiveWindow)
					// if name does not already exist in curSessionNamedWindow (like those transient windows we skipped earlier), include it.
					// The reason for this is because, this https://tronche.com/gui/x/icccm/sec-4.html#:~:text=It%20is%20important%20not,the%20window%20is%20mapped. might not be
					// adhered to by all applications. So, we are sure it can steal focus, so we include it.
					curSessionNamedWindow[currActiveWindow] = netActiveWindow.WindowName
				}

				if s.AppName != "" { // Mentioned earlier, if we don't have a name, lost metric.
					if err := x11.Db.WriteUsage(s); err != nil {
						log.Fatalf("write to db error:%v", err)
					}
				}
			}

		}
	}
}
