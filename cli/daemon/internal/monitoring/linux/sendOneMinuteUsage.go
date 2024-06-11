package monitoring

import (
	"context"
	"log"
	"pkg/types"
	"time"

	"github.com/BurntSushi/xgbutil/xevent"
)

func (x11 *X11Monitor) WindowChangeTimerFunc(ctx context.Context, timer *time.Timer) {
	for {
		select {
		case <-ctx.Done():
			return

		case <-x11.windowChangeCh:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(time.Duration(1) * time.Minute)

		case <-timer.C:
			timer.Reset(time.Duration(1) * time.Minute)
			x11.SendOneMinuteUsage()
		}
	}
}

func (x11 *X11Monitor) SendOneMinuteUsage() {

	if netActiveWindow.WindowID == xevent.NoWindow {
		return
	}

	oneMinuteUsage := time.Since(netActiveWindow.TimeStamp).Hours()
	oneMinuteTimeStamp := netActiveWindow.TimeStamp

	netActiveWindow.TimeStamp = time.Now()

	if err := x11.Db.WriteUsage(types.ScreenTime{
		WindowID: netActiveWindow.WindowID,
		AppName:  netActiveWindow.WindowName,
		Type:     types.Active,
		Duration: oneMinuteUsage,
		Interval: types.TimeInterval{Start: oneMinuteTimeStamp, End: time.Now()},
	}); err != nil {
		log.Fatalf("write to db error:%v", err)
	}
}

func (x11 *X11Monitor) CloseWindowChangeCh() {
	close(x11.windowChangeCh)
}
