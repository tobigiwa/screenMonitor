package database

import (
	"fmt"
	"utils"
)

var (
	ErrAppKeyMismatch = fmt.Errorf("key error: app name mismatch")
	ErrFutureWeek     = fmt.Errorf("requested date(week) is in the future, no time travelling with this func")
	ErrFutureDay      = fmt.Errorf("requested date(day) is in the future, no time travelling with this func")

	ZeroValueWeeklyStat = WeeklyStat{}
	ZeroValueDailyStat  = DailyStat{}
)

type dailyAppScreenTime map[utils.Date]utils.Stats

type AppInfo struct {
	utils.AppIconCategoryAndCmdLine
	ScreenStat dailyAppScreenTime `json:"screenStat"`
}

type DailyStat struct {
	EachApp  []utils.AppStat `json:"eachApp"`
	DayTotal utils.Stats     `json:"dayTotal"`
}

type WeeklyStat struct {
	EachApp       []utils.AppStat                                   `json:"eachApp"`
	WeekTotal     utils.Stats                                       `json:"weekTotal"`
	DayByDayTotal [7]utils.GenericKeyValue[utils.Date, utils.Stats] `json:"dayByDayTotal"`
}
