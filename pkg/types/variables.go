package types

import (
	"fmt"
	"regexp"
)

var (
	NoMessage                   = Message{}
	NoAppIconCategoryAndCmdLine = AppIconCategoryAndCmdLine{}

	InvalidDateType = Date("")
)

var (
	ErrDeserialization = fmt.Errorf("error deserializing data")
	ErrSerialization   = fmt.Errorf("error serializing data")
)

var (
	DateTypeRegexPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
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
