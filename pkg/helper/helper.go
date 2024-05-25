package helper

import (
	"encoding/json"
	"fmt"
	"math"
	"pkg/types"
	"strings"
	"time"
)

func Encode[T any](tyPe T) ([]byte, error) {
	encoded, err := json.Marshal(tyPe)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, types.ErrSerialization)
	}
	return encoded, nil
}

func Decode[T any](data []byte) (T, error) {
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
	return fmt.Sprintf("%.0f%%", percent(part, total))
}
