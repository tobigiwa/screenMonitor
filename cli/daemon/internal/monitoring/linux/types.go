package monitoring

import (
	db "LiScreMon/cli/daemon/internal/database"
	"sync"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
)

type netActiveWindowInfo struct {
	WindowID   xproto.Window
	WindowName string
	TimeStamp  time.Time
	DoNotCopy
}

type DoNotCopy [0]sync.Mutex

type X11Monitor struct {
	X11Connection *xgbutil.XUtil
	Db            db.IRepository
}
