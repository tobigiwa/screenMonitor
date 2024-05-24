package database

import (
	"fmt"
	"time"
)

func Key() Date {
	return Date(fmt.Sprint(time.Now().Format(timeFormat)))
}
func ParseKey(key Date) (time.Time, error) {
	return time.Parse(timeFormat, string(key))
}

var (
	dbAppPrefix  = []byte("app:")
	dbDayPrefix  = []byte("app:")
	dbWeekPrefix = []byte("week:")
)

func dbAppKey(appName string) []byte {
	return []byte(fmt.Sprintf("app:%v", appName))
}
func dbDayKey(date Date) []byte {
	return []byte(fmt.Sprintf("day:%v", string(date)))
}
func dbWeekKey(date Date) []byte {
	return []byte(fmt.Sprintf("week:%v", string(date)))
}
