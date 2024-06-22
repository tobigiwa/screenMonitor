package monitoring

import (
	"pkg/types"
	"sync"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/google/uuid"
)

type netActiveWindowInfo struct {
	windowInfo
	TimeStamp time.Time
	DoNotCopy
}

type DoNotCopy [0]sync.Mutex

type x11DBInterface interface {
	WriteUsage(types.ScreenTime) error
	UpdateOpertionOnBuCKET(dbPrefix string, opsFunc func([]byte) ([]byte, error)) error
	UpdateAppInfoManually(key []byte, opsFunc func([]byte) ([]byte, error)) error
	GetTaskByUUID(taskID uuid.UUID) (types.Task, error)
	RemoveTask(id uuid.UUID) error
	DeleteKey([]byte) error
	Close() error
}

type X11Monitor struct {
	windowChangeCh chan types.GenericKeyValue[xproto.Window, float64] //windowID and duration
	X11Connection  *xgbutil.XUtil
	Db             x11DBInterface
}

type windowInfo struct {
	WindowID   xproto.Window
	WindowName string
}

type limitWindow struct {
	windowInfo
	taskUUID       uuid.UUID
	timeSofar      float64
	limit          float64
	date           types.Date
	isLimitReached bool
}
