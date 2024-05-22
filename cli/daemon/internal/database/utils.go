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
)

func formattedToDay() time.Time {
	t, _ := ParseKey(Date(time.Now().Format(timeFormat)))
	return t
}

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

func SaturdayOfTheWeek(t time.Time) string {
	daysUntilSaturday := 6 - int(t.Weekday())
	return t.AddDate(0, 0, daysUntilSaturday).Format(timeFormat)
}

func daysInThatWeek(w time.Time) [7]Date {
	var arr [7]Date
	startOftheWeek := w.AddDate(0, 0, -int(w.Weekday()))
	for i := 0; i < 7; i++ {
		arr[i] = Date(fmt.Sprint(startOftheWeek.AddDate(0, 0, i).Format(timeFormat)))
	}
	return arr
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
