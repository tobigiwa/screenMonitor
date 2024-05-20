package service

import (
	db "LiScreMon/cli/daemon/internal/database"
	"fmt"
	"os"
	"pkg/types"
	"strings"
)

type Service struct {
	db db.IRepository
}

func (s *Service) getWeekStat(msg types.Message) types.WeekStatMessage {
	weekStat, err := s.db.GetWeek(msg.StringDataRequest)
	if err != nil {
		fmt.Println("error weekStat:", err)
		return types.WeekStatMessage{
			IsError: true,
			Error:   err,
		}
	}
	var (
		keys             = [7]string{}
		formattedDay     = [7]string{}
		values           = [7]float64{}
		sizeOfApps       = len(weekStat.EachApp)
		appNameInTheWeek = make([]string, 0, sizeOfApps)
		appCard          = make([]types.ApplicationDetail, 0, sizeOfApps)
	)

	for i := 0; i < 7; i++ {
		keys[i] = string(weekStat.DayByDayTotal[i].Key)
		values[i] = weekStat.DayByDayTotal[i].Value.Active

		day, _ := db.ParseKey(db.Date(weekStat.DayByDayTotal[i].Key))
		dayWithSuffix := addOrdinalSuffix(day.Day())
		weekDay := strings.TrimSuffix(day.Weekday().String(), "day")
		formattedDay[i] = fmt.Sprintf("%v. %v", weekDay, dayWithSuffix)
	}

	saturdayOftheWeek, _ := db.ParseKey(db.Date(weekStat.DayByDayTotal[6].Key))
	year, month, _ := saturdayOftheWeek.Date()
	stringMonth := month.String()

	for i := 0; i < sizeOfApps; i++ {
		appNameInTheWeek = append(appNameInTheWeek, weekStat.EachApp[i].AppName)
	}

	appsInfo, err := s.db.GetAppIconAndCategory(appNameInTheWeek)
	if err != nil {
		fmt.Println("err with GetAppIconAndCategory:", err)
		return types.WeekStatMessage{
			IsError: true,
			Error:   err,
		}
	}
	for i := 0; i < sizeOfApps; i++ {
		appCard = append(appCard, types.ApplicationDetail{AppInfo: appsInfo[i], Usage: weekStat.EachApp[i].Usage.Active})
	}

	return types.WeekStatMessage{
		Keys:            keys,
		FormattedDay:    formattedDay,
		Values:          values,
		TotalWeekUptime: weekStat.WeekTotal.Active,
		Month:           stringMonth,
		Year:            fmt.Sprint(year),
		AppDetail:       appCard,
	}
}

func addOrdinalSuffix(n int) string {
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

func savePNGImage(filename string, bytes []byte) error {
	// Create a new file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the byte slice to the file
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
