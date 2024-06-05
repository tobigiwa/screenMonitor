package service

import (
	db "LiScreMon/cli/daemon/internal/database"
	"fmt"
	helperFuncs "pkg/helper"
	"pkg/types"
)

type Service struct {
	db db.IRepository
}

func (s *Service) getWeekStat(msg types.Message) types.WeekStatMessage {
	var (
		weekStat db.WeeklyStat
		appsInfo []types.AppIconCategoryAndCmdLine
		err      error
	)

	if weekStat, err = s.db.GetWeek(msg.WeekStatRequest); err != nil {
		return types.WeekStatMessage{
			IsError: true,
			Error:   fmt.Errorf("error weekStat: %w", err)}
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
		formattedDay[i] = helperFuncs.FormattedDay(types.Date(weekStat.DayByDayTotal[i].Key))
	}

	month, year := helperFuncs.MonthAndYear(types.Date(weekStat.DayByDayTotal[6].Key))

	for i := 0; i < sizeOfApps; i++ {
		appNameInTheWeek = append(appNameInTheWeek, weekStat.EachApp[i].AppName)
	}

	if appsInfo, err = s.db.GetAppIconAndCategory(appNameInTheWeek); err != nil {
		return types.WeekStatMessage{
			IsError: true,
			Error:   fmt.Errorf("err with GetAppIconAndCategory:%w", err)}
	}

	for i := 0; i < sizeOfApps; i++ {
		appCard = append(appCard, types.ApplicationDetail{AppInfo: appsInfo[i], Usage: weekStat.EachApp[i].Usage.Active})
	}

	return types.WeekStatMessage{
		Keys:            keys,
		FormattedDay:    formattedDay,
		Values:          values,
		TotalWeekUptime: weekStat.WeekTotal.Active,
		Month:           month,
		Year:            fmt.Sprint(year),
		AppDetail:       appCard,
	}
}

func (s *Service) getAppStat(msg types.Message) types.AppStatMessage {
	var (
		appStat types.AppRangeStat
		err     error
	)

	switch msg.AppStatRequest.StatRange {
	case "week":
		appStat, err = s.db.AppWeeklyStat(msg.AppStatRequest.AppName, msg.AppStatRequest.Start)
	case "month":
		appStat, err = s.db.AppMonthlyStat(msg.AppStatRequest.AppName, msg.AppStatRequest.Month, msg.AppStatRequest.Year)
	case "range":
		appStat, err = s.db.AppDateRangeStat(msg.AppStatRequest.AppName, msg.AppStatRequest.Start, msg.AppStatRequest.End)
	}

	if err != nil {
		fmt.Println("error weekStat:", err)
		return types.AppStatMessage{
			IsError: true,
			Error:   err,
		}
	}

	var (
		formattedDay      = make([]string, 0, len(appStat.DaysRange))
		values            = []float64{}
		lastDayOfTheRange = len(appStat.DaysRange) - 1
	)

	for i := 0; i < len(appStat.DaysRange); i++ {
		formattedDay = append(formattedDay, helperFuncs.FormattedDay(types.Date(appStat.DaysRange[i].Key)))
		values = append(values, appStat.DaysRange[i].Value.Active)
	}
	month, year := helperFuncs.MonthAndYear(types.Date(appStat.DaysRange[lastDayOfTheRange].Key))

	return types.AppStatMessage{
		FormattedDay:     formattedDay,
		Values:           values,
		Month:            month,
		Year:             year,
		TotalRangeUptime: appStat.TotalRange.Active,
		AppInfo:          appStat.AppInfo,
	}
}

// func savePNGImage(filename string, bytes []byte) error {
// 	// Create a new file
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Write the byte slice to the file
// 	_, err = file.Write(bytes)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
