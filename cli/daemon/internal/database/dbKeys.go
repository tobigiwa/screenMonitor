package database

import (
	"fmt"
	"pkg/types"
	"time"
)

func today() types.Date {
	return types.Date(time.Now().Format(types.TimeFormat))
}
func yesterday() types.Date {
	return types.Date(time.Now().AddDate(0, 0, -1).Format(types.TimeFormat))
}


var (
	dbAppPrefix  = []byte("app:")
	dbDayPrefix  = []byte("day:")
	dbWeekPrefix = []byte("week:")
)

func dbAppKey(appName string) []byte {
	return []byte(fmt.Sprintf("app:%v", appName))
}
func dbDayKey(date types.Date) []byte {
	return []byte(fmt.Sprintf("day:%v", string(date)))
}
func dbWeekKey(date types.Date) []byte {
	return []byte(fmt.Sprintf("week:%v", string(date)))
}
func dbTaskKey() []byte {
	return []byte("tasks")
}
