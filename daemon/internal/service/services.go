package service

import (
	"LiScreMon/daemon/internal/database/repository"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strings"
)

type Message struct {
	Endpoint           string          `json:"endpoint"`
	StringDataRequest  string          `json:"stringDataRequest"`
	StringDataResponse string          `json:"stringDataResponse"`
	WeekStatResponse   WeekStatMessage `json:"weekStatResponse"`
}
type WeekStatMessage struct {
	Keys         [7]string  `json:"keys"`
	FormattedDay [7]string  `json:"formattedDay"`
	Values       [7]float64 `json:"values"`
	Month        string     `json:"month"`
	Year         string     `json:"year"`
	IsError      bool       `json:"isError"`
	Error        error      `json:"error"`
}

func (m *Message) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Message) decode(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := gob.NewDecoder(buf).Decode(m); err != nil {
		return err
	}
	return nil
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
		keys         = [7]string{}
		formattedDay = [7]string{}
		values       = [7]float64{}
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

	return WeekStatMessage{
		Keys:         keys,
		FormattedDay: formattedDay,
		Values:       values,
		Month:        stringMonth,
		Year:         fmt.Sprint(year),
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
