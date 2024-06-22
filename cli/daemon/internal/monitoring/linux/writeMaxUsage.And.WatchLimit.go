package monitoring

import (
	"context"
	"fmt"
	"log"
	"os"
	helperFuncs "pkg/helper"
	"pkg/types"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/gen2brain/beeep"
	"github.com/google/uuid"
)

func (x11 *X11Monitor) WindowChangeTimerFunc(ctx context.Context, timer *time.Timer) {
	for {
		select {
		case <-ctx.Done():
			return

		case t := <-x11.windowChangeCh:
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(fixtyEightSecs)
			x11.watchLimit(t.Key, t.Value)

		case <-timer.C:
			timer.Reset(fixtyEightSecs)

			if netActiveWindow.WindowID == xevent.NoWindow {
				return
			}
			x11.watchLimit(netActiveWindow.WindowID, fixtyEightSecs.Hours())
			x11.sendOneMinuteUsage()
		}
	}
}

func (x11 *X11Monitor) watchLimit(windowID xproto.Window, duration float64) {
	
	if windowName, ok := curSessionNamedWindow[windowID]; ok {
		if limitApp, ok := LimitApp[windowName]; ok {

			limitApp.timeSofar += duration

			if limitApp.timeSofar >= limitApp.limit {
				fmt.Printf("we have reached limit for this application\n%+v\n\n", limitApp)
				if err := x11.NotifyLimitReached(limitApp.taskUUID); err != nil {
					fmt.Println("error from notifyLimitReached", err)
				}

			} else {

				LimitApp[windowName] = limitApp
				fmt.Printf("\nthis so far %f for app %s...limit at %f\n\n", limitApp.timeSofar, windowName, limitApp.limit)
			}

			
		}
	}
}

func (x11 *X11Monitor) sendOneMinuteUsage() {

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

var LimitApp = make(map[string]limitWindow, 20)

func AddNewLimit(t types.Task) {
	LimitApp[t.AppName] = limitWindow{
		windowInfo: windowInfo{WindowName: t.AppName},
		taskUUID:   t.UUID,
		limit:      t.AppLimit.Limit,
	}
}

func (x11 *X11Monitor) CloseWindowChangeCh() {
	close(x11.windowChangeCh)
}

func (x11 *X11Monitor) NotifyLimitReached(taskID uuid.UUID) error {

	homeDir, _ := os.UserHomeDir()
	configDir := homeDir + "/liScreMon/"
	appLogo := configDir + "liscremon.jpeg"

	task, err := x11.Db.GetTaskByUUID(taskID)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("Usage Limit reached for %s", task.AppName)
	subtitle := fmt.Sprintf("App: %s Usage Limit: %s", task.AppName, helperFuncs.UsageTimeInHrsMin(task.AppLimit.Limit))

	beeep.Alert(title, subtitle, appLogo)

	if task.AppLimit.OneTime {
		if err := x11.Db.RemoveTask(taskID); err != nil {
			return err
		}
	}

	if err := x11.Db.UpdateAppLimitStatus(taskID); err != nil {
		return err
	}

	return nil
}
