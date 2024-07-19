package database

import (
	"fmt"
	"time"
	"utils"
)

func yesterday() utils.Date {
	return utils.Date(time.Now().AddDate(0, 0, -1).Format(utils.TimeFormat))
}

var (
	dbAppPrefix   = []byte("app:")
	dbDayPrefix   = []byte("day:")
	dbWeekPrefix  = []byte("week:")
	dbCategoryKey = []byte("category")
	dbTaskKey     = []byte("tasks")
)

func dbAppKey(appName string) []byte {
	return []byte(fmt.Sprintf("app:%v", appName))
}
func dbDayKey(date utils.Date) []byte {
	return []byte(fmt.Sprintf("day:%v", string(date)))
}
func dbWeekKey(date utils.Date) []byte {
	return []byte(fmt.Sprintf("week:%v", string(date)))
}
