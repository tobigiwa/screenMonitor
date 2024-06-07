package types

import "fmt"

var (
	NoMessage            = Message{}
	NoAppIconAndCategory = AppIconCategoryAndCmdLine{}
)

var (
	ErrDeserialization = fmt.Errorf("error deserializing data")
	ErrSerialization   = fmt.Errorf("error serializing data")
)

const (
	Active     ScreenType = "active"
	Inactive   ScreenType = "inactive"
	Open       ScreenType = "open"
	TimeFormat string     = "2006-01-02"

	ReminderWithNoAction TaskType = "ReminderWithNoAction"
	ReminderWithAction   TaskType = "ReminderWithAction"
	Limit                TaskType = "Limit"
)
