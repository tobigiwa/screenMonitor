package store

import (
	"fmt"
	"time"
)

func Key() date {
	now := time.Now()
	day, month, year := now.Day(), now.Month(), time.Now().Year()
	return date(fmt.Sprintf("%d:%d:%d", day, month, year))
}
func ParseKey(key string) time.Time {
	a, _ := time.Parse("2:1:2006", key)
	return a
}

func ShortenDay(t time.Weekday) string {

	switch t {
	case time.Sunday:
		return "Sun"
	case time.Monday:
		return "Mon"
	case time.Tuesday:
		return "Tue"
	case time.Wednesday:
		return "Wed"
	case time.Thursday:
		return "Thur"
	case time.Friday:
		return "Fri"
	case time.Saturday:
		return "Sat"
	}
	return ""
}
