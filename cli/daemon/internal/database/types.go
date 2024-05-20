package database

import (
	"fmt"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

var (
	ErrAppKeyMismatch = fmt.Errorf("key error: app name mismatch")
	ErrDeserilization = fmt.Errorf("error deserializing data")
	ErrSerilization   = fmt.Errorf("error serializing data")
	ErrFutureWeek     = fmt.Errorf("requested date(week) is in the future, no time travelling with this func")
	ErrFutureDay      = fmt.Errorf("requested date(day) is in the future, no time travelling with this func")
)

var (
	ZeroValueWeeklyStat = WeeklyStat{}
	ZeroValueDailyStat  = DailyStat{}
)

type ScreenType string

const (
	Active     ScreenType = "active"
	Inactive   ScreenType = "inactive"
	Open       ScreenType = "open"
	timeFormat string     = "2006-01-02"
)

type Category string

// date underneath is a
/* string of a time.Time format. "2006-01-02" */
type Date string
type dailyAppScreenTime map[Date]stats

type TimeInterval struct {
	Start time.Time
	End   time.Time
}
type stats struct {
	Active         float64
	Open           float64
	Inactive       float64
	ActiveTimeData []TimeInterval
}

// ScreenTime represents the time spent on a particular app.
type ScreenTime struct {
	WindowID xproto.Window
	AppName  string
	Type     ScreenType
	Duration float64
	Interval TimeInterval
}

type appStat struct {
	AppName string
	Usage   stats
}

type DailyStat struct {
	EachApp  []appStat
	DayTotal stats
}

type GenericKeyValue[K, V any] struct {
	Key   K
	Value V
}

type WeeklyStat struct {
	EachApp       []appStat
	WeekTotal     stats
	DayByDayTotal [7]GenericKeyValue[Date, stats]
}

type appInfo struct {
	AppName           string
	Icon              []byte
	IsIconSet         bool
	Category          Category
	IsCategorySet     bool
	DesktopCategories []string
	ScreenStat        dailyAppScreenTime
}
