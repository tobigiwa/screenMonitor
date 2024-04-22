package store

import (
	"bytes"
	"encoding/gob"
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

type ScreenType string

const (
	Active  ScreenType = "active"
	Passive ScreenType = "passive"

	timeFormat string = "2006-01-02"
)

// ScreenTime represents the time spent on a particular app.
type ScreenTime struct {
	WindowID xproto.Window
	AppName  string
	Type     ScreenType
	Duration float64
	Interval TimeInterval
}

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

type Category string

type appInfo struct {
	AppName           string
	Icon              []byte
	IsIconSet         bool
	Category          Category
	IsCategorySet     bool
	DesktopCategories []string
	ScreenStat        dailyAppScreenTime
}

func (ap *appInfo) serialize() ([]byte, error) {
	var res bytes.Buffer
	encoded := gob.NewEncoder(&res)
	if err := encoded.Encode(ap); err != nil {
		return nil, fmt.Errorf("%v:%w", err, ErrSerilization)
	}
	return res.Bytes(), nil
}

func (ap *appInfo) deserialize(data []byte) error {
	decoded := gob.NewDecoder(bytes.NewReader(data))
	if err := decoded.Decode(ap); err != nil {
		return fmt.Errorf("%v:%w", err, ErrDeserilization)
	}
	return nil
}

func Key() date {
	return date(fmt.Sprint(time.Now().Format(timeFormat)))
}
func ParseKey(key date) time.Time {
	a, _ := time.Parse(timeFormat, string(key))
	return a
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
