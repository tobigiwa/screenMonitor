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
type dailyAppScreenTime map[Date]Stats

type AppInfo struct {
	AppName           string             `json:"appName"`
	Icon              []byte             `json:"icon"`
	IsIconSet         bool               `json:"isIconSet"`
	Category          Category           `json:"category"`
	IsCategorySet     bool               `json:"isCategorySet"`
	DesktopCategories []string           `json:"desktopCategories"`
	ScreenStat        dailyAppScreenTime `json:"screenStat"`
}

type TimeInterval struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type ScreenTime struct {
	WindowID xproto.Window `json:"windowID"`
	AppName  string        `json:"appName"`
	Type     ScreenType    `json:"type"`
	Duration float64       `json:"duration"`
	Interval TimeInterval  `json:"interval"`
}

type Stats struct {
	Active         float64        `json:"active"`
	Open           float64        `json:"open"`
	Inactive       float64        `json:"inactive"`
	ActiveTimeData []TimeInterval `json:"activeTimeData"`
}

type AppStat struct {
	AppName string `json:"appName"`
	Usage   Stats  `json:"usage"`
}

type GenericKeyValue[K, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}

type DailyStat struct {
	EachApp  []AppStat `json:"eachApp"`
	DayTotal Stats     `json:"dayTotal"`
}

type WeeklyStat struct {
	EachApp       []AppStat                       `json:"eachApp"`
	WeekTotal     Stats                           `json:"weekTotal"`
	DayByDayTotal [7]GenericKeyValue[Date, Stats] `json:"dayByDayTotal"`
}
