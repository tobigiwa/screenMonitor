package database

import (
	"fmt"
	"time"
)

func Key() Date {
	return Date(fmt.Sprint(time.Now().Format(timeFormat)))
}
func ParseKey(key Date) (time.Time, error) {
	a, err := time.Parse(timeFormat, string(key))
	if err != nil {
		return time.Time{}, err
	}
	return a, nil
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
