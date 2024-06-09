package monitoring

import (
	"log"
	"pkg/types"
	"time"

	"github.com/BurntSushi/xgbutil/xevent"
)

func (x11 *X11Monitor) windowChangeTimerFunc() {
	defer func() {
		if !x11.timer.Stop() {
			<-x11.timer.C
		}
	}()

	for {
		select {
		case <-x11.ctx.Done():
			return

		case <-x11.windowChangeCh:
			if !x11.timer.Stop() {
				<-x11.timer.C
			}
			x11.timer.Reset(time.Duration(1) * time.Minute)

		case <-x11.timer.C:
			x11.timer.Reset(time.Duration(1) * time.Minute)
			x11.SendOneMinuteUsage()
		}
	}
}

func (x11 *X11Monitor) SendOneMinuteUsage() {

	if netActiveWindow.WindowID == xevent.NoWindow {
		return
	}

	timeSoFar := time.Since(netActiveWindow.TimeStamp).Hours()
	timeStartTimeStamp := netActiveWindow.TimeStamp

	netActiveWindow.TimeStamp = time.Now()

	if err := x11.Db.WriteUsage(types.ScreenTime{
		WindowID: netActiveWindow.WindowID,
		AppName:  netActiveWindow.WindowName,
		Type:     types.Active,
		Duration: timeSoFar,
		Interval: types.TimeInterval{Start: timeStartTimeStamp, End: time.Now()},
	}); err != nil {
		log.Fatalf("write to db error:%v", err)
	}
}

func (x11 *X11Monitor) CloseWindowChangeCh() {
	close(x11.windowChangeCh)
}
