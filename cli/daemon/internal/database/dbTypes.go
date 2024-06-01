package database

import (
	"fmt"

	"pkg/types"
)

var (
	ErrAppKeyMismatch = fmt.Errorf("key error: app name mismatch")
	ErrFutureWeek     = fmt.Errorf("requested date(week) is in the future, no time travelling with this func")
	ErrFutureDay      = fmt.Errorf("requested date(day) is in the future, no time travelling with this func")

	ZeroValueWeeklyStat = WeeklyStat{}
	ZeroValueDailyStat  = DailyStat{}
)

type dailyAppScreenTime map[types.Date]types.Stats

type AppInfo struct {
	AppName           string             `json:"appName"`
	IsIconSet         bool               `json:"isIconSet"`
	IsCmdLineSet      bool               `json:"isCmdLine`
	IsCategorySet     bool               `json:"isCategorySet"`
	cmdLine           string             `json:"cmdLine`
	Icon              []byte             `json:"icon"`
	DesktopCategories []string           `json:"desktopCategories"`
	Category          types.Category     `json:"category"`
	ScreenStat        dailyAppScreenTime `json:"screenStat"`
}

type AppStat struct {
	AppName string      `json:"appName"`
	Usage   types.Stats `json:"usage"`
}

type DailyStat struct {
	EachApp  []AppStat   `json:"eachApp"`
	DayTotal types.Stats `json:"dayTotal"`
}

type WeeklyStat struct {
	EachApp       []AppStat                                         `json:"eachApp"`
	WeekTotal     types.Stats                                       `json:"weekTotal"`
	DayByDayTotal [7]types.GenericKeyValue[types.Date, types.Stats] `json:"dayByDayTotal"`
}
