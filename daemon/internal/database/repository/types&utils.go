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
type date string
type dailyAppScreenTime map[date]stats

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

func Key() date {
	return date(fmt.Sprint(time.Now().Format(timeFormat)))
}
func ParseKey(key date) (time.Time, error) {
	a, err := time.Parse(timeFormat, string(key))
	if err != nil {
		return time.Time{}, err
	}
	return a, nil
}

func dbAppKey(appName string) []byte {
	return []byte(fmt.Sprintf("app:%v", appName))
}

func dbAppPrefix() []byte {
	return []byte("app:")
}

func dbActiveSTKey() []byte {
	return []byte("active")
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

func availableStatForThatWeek(w time.Time) []date {
	var n int
	switch w.Weekday() {
	case time.Monday:
		n = 1
	case time.Tuesday:
		n = 2
	case time.Wednesday:
		n = 3
	case time.Thursday:
		n = 4
	case time.Friday:
		n = 5
	case time.Saturday:
		n = 6
	case time.Sunday:
		n = 0
	}
	arr := make([]date, 0, n)
	for i := 0; i <= n; i++ {
		arr = append(arr, date(fmt.Sprint(w.AddDate(0, 0, -i).Format(timeFormat))))
	}
	return arr
}
