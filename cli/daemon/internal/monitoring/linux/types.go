package monitoring

import (
	"pkg/types"
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

type x11DBInterface interface {
	WriteUsage(types.ScreenTime) error
	UpdateOpertionOnBuCKET(dbPrefix string, opsFunc func([]byte) ([]byte, error)) error
	DeleteKey([]byte) error
	Close() error
}

type X11Monitor struct {
	windowChangeCh chan struct{}
	X11Connection  *xgbutil.XUtil
	Db             x11DBInterface
}
