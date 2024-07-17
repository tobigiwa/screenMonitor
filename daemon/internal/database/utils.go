package database

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
	"utils"
)

func IsFutureWeek(t time.Time) bool {
	now := time.Now()

	// Adjust the week start day to Sunday
	now = now.AddDate(0, 0, int(time.Sunday-now.Weekday()))
	t = t.AddDate(0, 0, int(time.Sunday-t.Weekday()))

	_, currentWeek := now.ISOWeek()
	_, inputWeek := t.ISOWeek()

	return inputWeek > currentWeek
}

func IsPastWeek(t time.Time) bool {
	now := time.Now()

	// Adjust the week start day to Sunday
	now = now.AddDate(0, 0, int(time.Sunday-now.Weekday()))
	t = t.AddDate(0, 0, int(time.Sunday-t.Weekday()))

	_, currentWeek := now.ISOWeek()
	_, inputWeek := t.ISOWeek()

	return inputWeek < currentWeek
}

func daysInThatWeek(w time.Time) [7]utils.Date {
	var arr [7]utils.Date
	startOftheWeek := w.AddDate(0, 0, -int(w.Weekday()))
	for i := 0; i < 7; i++ {
		arr[i] = utils.Date(fmt.Sprint(startOftheWeek.AddDate(0, 0, i).Format(utils.TimeFormat)))
	}
	return arr
}

func getDesktopCategoryAndCmd(appName string) (dotDesktopFileInfo, error) {
	var r dotDesktopFileInfo

	if OperatingSytem := runtime.GOOS; OperatingSytem == "linux" {
		dir := "/usr/share/applications/"
		files, err := os.ReadDir(dir)
		if err != nil {
			return dotDesktopFileInfo{}, err
		}

		for _, file := range files {
			if strings.Contains(strings.ToLower(file.Name()), strings.ToLower(appName)) && strings.HasSuffix(file.Name(), ".desktop") {
				content, err := os.ReadFile(filepath.Join(dir, file.Name()))
				if err != nil {
					// continue
					// since there should be only one
					return dotDesktopFileInfo{}, err
				}

				lines := bytes.Split(content, []byte("\n"))
				for i := 0; i < len(lines); i++ {
					line := string(lines[i])

					if strings.HasPrefix(line, "Exec=") {
						r.cmdLine = strings.TrimPrefix(line, "Exec=")
					}

					if strings.HasPrefix(line, "Categories=") {
						if after, found := strings.CutPrefix(line, "Categories="); found {
							categories := strings.Split(after, ";")

							// trims out empty value, some end the line with ";"
							categories = slices.DeleteFunc(categories, func(s string) bool {
								return strings.TrimSpace(s) == ""
							})

							r.desktopCategories = categories
						}
					}
					if r.cmdLine != "" && r.desktopCategories != nil {
						return r, nil
					}

				}
				// since there should be only one .desktop for a name
				return r, nil // return anyone of 'em that has been set
			}
		}

	} else if OperatingSytem == "windows" {
		return dotDesktopFileInfo{}, nil
	}

	return dotDesktopFileInfo{}, errors.New("just an error")
}

type dotDesktopFileInfo struct {
	desktopCategories []string
	cmdLine           string
}
