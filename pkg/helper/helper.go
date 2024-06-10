package helper

import (
	"encoding/json"
	"fmt"
	"math"
	"pkg/types"
	"strings"
	"time"
)

func EncodeJSON[T any](tyPe T) ([]byte, error) {
	encoded, err := json.Marshal(tyPe)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, types.ErrSerialization)
	}
	return encoded, nil
}

func DecodeJSON[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("%v:%w", err, types.ErrDeserialization)
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

func ParseKey(key types.Date) (time.Time, error) {
	return time.Parse(types.TimeFormat, string(key))
}

func FormattedDay(date types.Date) string {
	day, _ := ParseKey(date)
	dayWithSuffix := AddOrdinalSuffix(day.Day())
	weekDay := strings.TrimSuffix(day.Weekday().String(), "day")
	return fmt.Sprintf("%v. %v", weekDay, dayWithSuffix)
}

func MonthAndYear(date types.Date) (month, year string) {
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

func SaturdayOfTheWeek(t time.Time) string {
	daysUntilSaturday := 6 - int(t.Weekday())
	return t.AddDate(0, 0, daysUntilSaturday).Format(types.TimeFormat)
}

func FirstSaturdayOfTheMonth(month string) string {
	t, err := time.Parse("January", month)
	if err != nil {
		return ""
	}
	year := time.Now().Year()
	firstDayOfMonth := time.Date(year, t.Month(), 1, 0, 0, 0, 0, time.UTC)

	for {
		if firstDayOfMonth.Weekday() == time.Saturday {
			return firstDayOfMonth.Format(types.TimeFormat)
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
	return s.Format(types.TimeFormat)
}

func ReturnLastWeekSaturday(t time.Time) string {

	if t.Weekday() == time.Saturday {
		return t.AddDate(0, 0, -7).Format(types.TimeFormat)
	}

	daysSinceSaturday := int(t.Weekday()+1) % 7
	return t.AddDate(0, 0, -daysSinceSaturday).Format(types.TimeFormat)
}

func IsFutureDate(t time.Time) bool {
	today := time.Now()
	nextWeekDay := t.AddDate(0, 0, 7)
	return nextWeekDay.After(today)
}

func ReturnNexWeektSaturday(saturday time.Time) string {
	return saturday.AddDate(0, 0, 7).Format(types.TimeFormat)
}

func FormattedToDay() time.Time {
	t, _ := ParseKey(types.Date(time.Now().Format(types.TimeFormat)))
	return t
}

func AllTheDaysInMonth(year, month string) ([]types.Date, error) {
	t, err := time.Parse("2006 January", year+" "+month)
	if err != nil {
		return nil, fmt.Errorf("parse %w", err)
	}

	fmt.Println(t.Day(), t.Month(), t.Year())
	lastDayOfTheGivenMonth := time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()

	dates := make([]types.Date, 0, lastDayOfTheGivenMonth)

	for day := 1; day <= lastDayOfTheGivenMonth; day++ {
		dates = append(dates, types.Date(time.Date(t.Year(), t.Month(), day, 0, 0, 0, 0, t.Location()).Format(types.TimeFormat)))
	}

	return dates, nil
}

func ValidDateType(s string) bool {
	return types.DateTypeRegexPattern.MatchString(s)
}
