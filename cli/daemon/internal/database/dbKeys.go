package database

import (
	"fmt"
	"pkg/types"
	"time"
)

func Key() types.Date {
	return types.Date(fmt.Sprint(time.Now().Format(types.TimeFormat)))
}
func ParseKey(key types.Date) (time.Time, error) {
	return time.Parse(types.TimeFormat, string(key))
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
