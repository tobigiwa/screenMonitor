package service

import (
	"LiScreMon/daemon/internal/database/repository"
	"fmt"
	"log"
	"os"
	"strings"
)

type Message struct {
	Endpoint           string          `json:"endpoint"`
	StringDataRequest  string          `json:"stringDataRequest"`
	StringDataResponse string          `json:"stringDataResponse"`
	WeekStatResponse   WeekStatMessage `json:"weekStatResponse"`
}
type WeekStatMessage struct {
	Keys            [7]string           `json:"keys"`
	FormattedDay    [7]string           `json:"formattedDay"`
	Values          [7]float64          `json:"values"`
	TotalWeekUptime float64             `json:"totalWeekUptime"`
	Month           string              `json:"month"`
	Year            string              `json:"year"`
	AppDetail       []applicationDetail `json:"appDetail"`
	IsError         bool                `json:"isError"`
	Error           error               `json:"error"`
}

type applicationDetail struct {
	AppInfo repository.AppIconAndCategory `json:"appInfo"`
	Usage   float64                       `json:"usage"`
}

type Service struct {
	store repository.IRepository
}

func (s *Service) getWeekStat(msg Message) WeekStatMessage {
	weekStat, err := s.store.GetWeek(msg.StringDataRequest)
	if err != nil {
		log.Println("error weekStat:", err)
		return WeekStatMessage{
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
		appCard          = make([]applicationDetail, 0, sizeOfApps)
	)

	for i := 0; i < 7; i++ {
		keys[i] = string(weekStat.DayByDayTotal[i].Key)
		values[i] = weekStat.DayByDayTotal[i].Value.Active

		day, _ := repository.ParseKey(repository.Date(weekStat.DayByDayTotal[i].Key))
		dayWithSuffix := addOrdinalSuffix(day.Day())
		weekDay := strings.TrimSuffix(day.Weekday().String(), "day")
		formattedDay[i] = fmt.Sprintf("%v. %v", weekDay, dayWithSuffix)
	}

	saturdayOftheWeek, _ := repository.ParseKey(repository.Date(weekStat.DayByDayTotal[6].Key))
	year, month, _ := saturdayOftheWeek.Date()
	stringMonth := month.String()

	for i := 0; i < sizeOfApps; i++ {
		appNameInTheWeek = append(appNameInTheWeek, weekStat.EachApp[i].AppName)
	}

	appsInfo, err := s.store.GetAppIconAndCategory(appNameInTheWeek)
	if err != nil {
		fmt.Println("err with GetAppIconAndCategory:", err)
	}
	for i := 0; i < sizeOfApps; i++ {
		appCard = append(appCard, applicationDetail{AppInfo: appsInfo[i], Usage: weekStat.EachApp[i].Usage.Active})
	}

	return WeekStatMessage{
		Keys:            keys,
		FormattedDay:    formattedDay,
		Values:          values,
		TotalWeekUptime: weekStat.WeekTotal.Active,
		Month:           stringMonth,
		Year:            fmt.Sprint(year),
		// AppDetail:       appCard,
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
