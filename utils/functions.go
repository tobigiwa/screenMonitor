package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

func EncodeJSON[T any](tyPe T) ([]byte, error) {
	encoded, err := json.Marshal(tyPe)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, ErrSerialization)
	}
	return encoded, nil
}

func DecodeJSON[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("%v:%w", err, ErrDeserialization)
	}
	return result, nil
}

func HrsAndMinute(hr float64) (int, int) {
	return int(hr), int(math.Round((hr - float64(int(hr))) * 60))
}

func UsageTimeInHrsMin(f float64) string {
	hrs, min := HrsAndMinute(f)
	return fmt.Sprintf("%dHrs:%dMin", hrs, min)
}

func AddOrdinalSuffix(n int) string {
	switch n {
	case 1, 21, 31:
		return fmt.Sprintf("%dst", n)
	case 2, 22:
		return fmt.Sprintf("%dnd", n)
	case 3, 23:
		return fmt.Sprintf("%drd", n)
	default:
		return fmt.Sprintf("%dth", n)
	}
}

func ParseKey(key Date) (time.Time, error) {
	return time.Parse(TimeFormat, string(key))
}
func Today() Date {
	return Date(time.Now().Format(TimeFormat))
}
func FormattedDay(date Date) string {
	day, _ := ParseKey(date)
	dayWithSuffix := AddOrdinalSuffix(day.Day())
	weekDay := strings.TrimSuffix(day.Weekday().String(), "day")
	return fmt.Sprintf("%v. %v", weekDay, dayWithSuffix)
}

func MonthAndYear(date Date) (month, year string) {
	day, _ := ParseKey(date)
	y, mon, _ := day.Date()
	month, year = mon.String(), fmt.Sprint(y)
	return month, year
}

func percent(part, total float64) float64 {
	return (part / total) * 100
}

func PercentagesString(part, total float64) string {
	p := percent(part, total)
	if math.IsNaN(p) {
		return "NaN"
	}
	return fmt.Sprintf("%.0f%%", p)
}

func SaturdayOfTheWeek(t time.Time) Date {
	daysUntilSaturday := 6 - int(t.Weekday())
	return Date(t.AddDate(0, 0, daysUntilSaturday).Format(TimeFormat))
}

func FirstSaturdayOfTheMonth(month string) Date {
	t, err := time.Parse("January", month)
	if err != nil {
		return ""
	}
	year := time.Now().Year()
	firstDayOfMonth := time.Date(year, t.Month(), 1, 0, 0, 0, 0, time.UTC)

	for {
		if firstDayOfMonth.Weekday() == time.Saturday {
			return Date(firstDayOfMonth.Format(TimeFormat))
		}
		firstDayOfMonth = firstDayOfMonth.AddDate(0, 0, 1)
	}
}

func LastSaturdayOfTheMonth(month string) string {
	t, err := time.Parse("January", month)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	NextMonth := time.Date(time.Now().Year(), t.Month()+1, 1, 0, 0, 0, 0, time.UTC)

	var s time.Time
	for {
		NextMonth = NextMonth.AddDate(0, 0, -1)
		if NextMonth.Weekday() == time.Saturday {
			s = NextMonth
			break
		}
	}
	return s.Format(TimeFormat)
}

func ReturnLastWeekSaturday(t time.Time) Date {

	if t.Weekday() == time.Saturday {
		return Date(t.AddDate(0, 0, -7).Format(TimeFormat))
	}

	daysSinceSaturday := int(t.Weekday()+1) % 7
	return Date(t.AddDate(0, 0, -daysSinceSaturday).Format(TimeFormat))
}

func IsFutureDate(t time.Time) bool {
	today := time.Now()
	nextWeekDay := t.AddDate(0, 0, 7)
	return nextWeekDay.After(today)
}

func ReturnNexWeektSaturday(saturday time.Time) Date {
	return Date(saturday.AddDate(0, 0, 7).Format(TimeFormat))
}

func FormattedToDay() time.Time {
	t, _ := ParseKey(Date(time.Now().Format(TimeFormat)))
	return t
}

func AllTheDaysInMonth(year, month string) ([]Date, error) {
	t, err := time.Parse("2006 January", year+" "+month)
	if err != nil {
		return nil, fmt.Errorf("parse %w", err)
	}

	fmt.Println(t.Day(), t.Month(), t.Year())
	lastDayOfTheGivenMonth := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()

	dates := make([]Date, 0, lastDayOfTheGivenMonth)

	for day := 1; day <= lastDayOfTheGivenMonth; day++ {
		dates = append(dates, Date(time.Date(t.Year(), t.Month(), day, 0, 0, 0, 0, t.Location()).Format(TimeFormat)))
	}

	return dates, nil
}

func ValidDateType(s string) bool {
	return DateTypeRegexPattern.MatchString(s)
}

func IsInCurrentWeekTime(t time.Time) bool {
	now := time.Now()
	_, currentWeek := now.ISOWeek()
	_, tWeek := t.ISOWeek()

	return currentWeek == tWeek && now.Year() == t.Year()
}

func IsInCurrentWeekDate(d Date) bool {
	// t, _ := d.ToTime()
	return SaturdayOfTheWeek(time.Now()) == d
}

func ConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, "liScreMon")

	AppLogoFilePath = filepath.Join(configDir, "liscremon.jpeg")

	return configDir, nil
}

func JSONConfigFile() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

func NotifyWithBeep(title, subtitle string) {
	beeep.Alert(title, subtitle, AppLogoFilePath)
}
func NotifyWithoutBeep(title, subtitle string) {
	beeep.Notify(title, subtitle, AppLogoFilePath)
}
