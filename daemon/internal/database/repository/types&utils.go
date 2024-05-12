package repository

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/BurntSushi/xgb/xproto"
)

var (
	ErrAppKeyMismatch = fmt.Errorf("key error: app name mismatch")
	ErrDeserilization = fmt.Errorf("error deserializing data")
	ErrSerilization   = fmt.Errorf("error serializing data")

	ErrKeyNotFound = fmt.Errorf("data does not exist")

	ErrFutureWeek = fmt.Errorf("requested date(week) is in the future, no time travelling with this func")
	ErrFutureDay = fmt.Errorf("requested date(day) is in the future, no time travelling with this func")
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
type dailyAppScreenTime map[Date]stats

type TimeInterval struct {
	Start time.Time
	End   time.Time
}
type stats struct {
	Active         float64
	Open           float64
	Inactive       float64
	ActiveTimeData []TimeInterval
}
type dailyActiveScreentime struct {
	Stats stats
}

// ScreenTime represents the time spent on a particular app.
type ScreenTime struct {
	WindowID xproto.Window
	AppName  string
	Type     ScreenType
	Duration float64
	Interval TimeInterval
}

type KeyValuePair struct {
	Key   string
	Value float64
}

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

var dbAppPrefix = []byte("app:")
var dbDayPrefix = []byte("app:")
var dbWeekPrefix = []byte("week:")


func dbAppKey(appName string) []byte {
	return []byte(fmt.Sprintf("app:%v", appName))
}
func dbDayKey(date Date) []byte {
	return []byte(fmt.Sprintf("day:%v", string(date)))
}
func dbWeekKey(date Date) []byte {
	return []byte(fmt.Sprintf("week:%v", string(date)))
}

func getDesktopCategory(appName string) ([]string, error) {

	if OperatingSytem := runtime.GOOS; OperatingSytem == "linux" {
		dir := "/usr/share/applications/"
		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(appName)) && strings.HasSuffix(file.Name(), ".desktop") {
				content, err := os.ReadFile(filepath.Join(dir, file.Name()))
				if err != nil {
					continue
				}
				lines := bytes.Split(content, []byte("\n"))
				for i := 0; i < len(lines); i++ {
					if line := string(lines[i]); strings.HasPrefix(line, "Categories=") {
						if after, found := strings.CutPrefix(line, "Categories="); found {
							categories := strings.Split(after, ";")

							categories = slices.DeleteFunc(categories, func(s string) bool { // some end the line with ";"
								return strings.TrimSpace(s) == ""
							})
							return categories, nil
						}
					}
				}
			}
		}

	} else if OperatingSytem == "windows" {
		return nil, nil
	}

	return nil, errors.New("just an error")
}

func daysInThatWeek(w time.Time) [7]Date {

	var arr [7]Date
	startOftheWeek := w.AddDate(0, 0, -int(w.Weekday()))
	for i := 0; i < 7; i++ {
		arr[i] = Date(fmt.Sprint(startOftheWeek.AddDate(0, 0, i).Format(timeFormat)))
	}
	return arr
}
