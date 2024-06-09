package monitoring

import (
	"context"
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
	Close() error
}

type X11Monitor struct {
	ctx            context.Context
	timer          *time.Timer
	windowChangeCh chan struct{}
	CancelFunc     context.CancelFunc
	X11Connection  *xgbutil.XUtil
	Db             x11DBInterface
}
