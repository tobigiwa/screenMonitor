package types

import (
	"regexp"
)

var (
	NoMessage                   = Message{}
	NoAppIconCategoryAndCmdLine = AppIconCategoryAndCmdLine{}

	InvalidDateType = Date("")
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
	DailyAppLimit        TaskType = "DailyAppLimit"
)

const (
	ProductivityAndUtility           Category = "Productivity & Utility"
	CommunicationAndSocialNetworking Category = "Communication & Social Networking"
	EntertainmentAndGaming           Category = "Entertainment & Gaming"
	WebBrowser                       Category = "Web Browser"
	SytemApp                         Category = "System App"
)

var CategoryMap map[string]Category = map[string]Category{
	"utilities":        ProductivityAndUtility,
	"productivity":     ProductivityAndUtility,
	"texteditor":       ProductivityAndUtility,
	"development":      ProductivityAndUtility,
	"ide":              ProductivityAndUtility,
	"editor":           ProductivityAndUtility,
	"viewer":           ProductivityAndUtility,
	"office":           ProductivityAndUtility,
	"communication":    CommunicationAndSocialNetworking,
	"social":           CommunicationAndSocialNetworking,
	"instantmessaging": CommunicationAndSocialNetworking,
	"entertainment":    EntertainmentAndGaming,
	"gaming":           EntertainmentAndGaming,
	"player":           EntertainmentAndGaming,
	"video":            EntertainmentAndGaming,
	"audio":            EntertainmentAndGaming,
	"audiovideo":       EntertainmentAndGaming,
	"browser":          WebBrowser,
	"webbrowser":       WebBrowser,
	"gnome":            SytemApp,
	"terminal":         SytemApp,
	"system":           SytemApp,
	"filemanager":      SytemApp,
	"filesystem":       SytemApp,
	"core":             SytemApp,
	"packagemanager":   SytemApp,
	"settings":         SytemApp,
	"terminalemulator": SytemApp,
}

var DefalutCategory = []Category{ProductivityAndUtility, CommunicationAndSocialNetworking,
	EntertainmentAndGaming, WebBrowser, SytemApp}
