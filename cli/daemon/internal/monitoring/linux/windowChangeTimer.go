package monitoring

import (
	"time"
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

			x11.timer.Reset(time.Duration(10) * time.Second)

		case <-x11.timer.C:
			sendOneMinuteUsage()
			x11.timer.Reset(time.Duration(10) * time.Second)
		}
	}
}

func sendOneMinuteUsage() {
}

func (x11 *X11Monitor) CloseWindowChangeCh() {
	close(x11.windowChangeCh)
}
