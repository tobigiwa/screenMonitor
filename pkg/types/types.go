package types

import (
	"bytes"
	"encoding/binary"
)

type Message struct {
	Endpoint           string
	StringDataRequest  string
	StringDataResponse string
	WeekStatResponse   WeekStatMessage
}
type WeekStatMessage struct {
	Keys            [7]string
	FormattedDay    [7]string
	Values          [7]float64
	TotalWeekUptime float64
	Month           string
	Year            string
	AppDetail       []ApplicationDetail
	IsError         bool
	Error           error
}

type AppIconAndCategory struct {
	AppName           string
	Icon              []byte
	IsIconSet         bool
	Category          string
	IsCategorySet     bool
	DesktopCategories []string
}

type ApplicationDetail struct {
	AppInfo AppIconAndCategory
	Usage   float64
}

func migrateOldData(oldBytes []byte) (Message, error) {
	var oldData Message
	if err := binary.Read(bytes.NewReader(oldBytes), binary.LittleEndian, &oldData); err != nil {
		return Message{}, err
	}

	// Map fields to the new struct
	newData := Message{
		Endpoint:           oldData.Endpoint,
		StringDataRequest:  oldData.StringDataRequest,
		StringDataResponse: oldData.StringDataResponse,
		WeekStatResponse:   oldData.WeekStatResponse,
		// Initialize other fields...
	}

	return newData, nil
}
