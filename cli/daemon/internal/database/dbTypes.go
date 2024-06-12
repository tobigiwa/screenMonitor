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
	types.AppIconCategoryAndCmdLine
	ScreenStat dailyAppScreenTime `json:"screenStat"`
}

type DailyStat struct {
	EachApp  []types.AppStat `json:"eachApp"`
	DayTotal types.Stats     `json:"dayTotal"`
}

type WeeklyStat struct {
	EachApp       []types.AppStat                                   `json:"eachApp"`
	WeekTotal     types.Stats                                       `json:"weekTotal"`
	DayByDayTotal [7]types.GenericKeyValue[types.Date, types.Stats] `json:"dayByDayTotal"`
}
