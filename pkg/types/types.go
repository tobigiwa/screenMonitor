package types

import (
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/google/uuid"
)

type Message struct {
	Endpoint         string          `json:"endpoint"`
	StatusCheck      string          `json:"statusCheck"`
	WeekStatRequest  string          `json:"weekStatRequest"`
	WeekStatResponse WeekStatMessage `json:"weekStatResponse"`
	AppStatRequest   AppStatRequest  `json:"appStatResquest"`
	AppStatResponse  AppStatMessage  `json:"appStatResponse"`
}
type WeekStatMessage struct {
	Keys            [7]string           `json:"keys"`
	FormattedDay    [7]string           `json:"formattedDay"`
	Values          [7]float64          `json:"values"`
	TotalWeekUptime float64             `json:"totalWeekUptime"`
	Month           string              `json:"month"`
	Year            string              `json:"year"`
	AppDetail       []ApplicationDetail `json:"appDetail"`
	IsError         bool                `json:"isError"`
	Error           error               `json:"error"`
}

type AppStatRequest struct {
	AppName   string `json:"appName"`
	Month     string `json:"month"`
	Year      string `json:"year"`
	StatRange string `json:"statRange"`
	Start     Date   `json:"start"`
	End       Date   `json:"end"`
}

type AppStatMessage struct {
	FormattedDay     []string           `json:"formattedDay"`
	Values           []float64          `json:"values"`
	Month            string             `json:"month"`
	Year             string             `json:"year"`
	TotalRangeUptime float64            `json:"totalRangeUptime"`
	AppInfo          AppIconAndCategory `json:"appInfo"`
	IsError          bool               `json:"isError"`
	Error            error              `json:"error"`
}

type AppIconAndCategory struct {
	AppName           string   `json:"appName"`
	Icon              []byte   `json:"icon"`
	IsIconSet         bool     `json:"isIconSet"`
	Category          string   `json:"category"`
	IsCategorySet     bool     `json:"isCategorySet"`
	DesktopCategories []string `json:"desktopCategories"`
}

type ApplicationDetail struct {
	AppInfo      AppIconAndCategory `json:"appInfo"`
	Usage        float64            `json:"usage"`
	AnyDayInStat Date               `json:"anyDayInStat"`
}

type GenericKeyValue[K, V any] struct {
	Key   K `json:"key"`
	Value V `json:"value"`
}

type AppRangeStat struct {
	AppInfo    AppIconAndCategory             `json:"appInfo"`
	DaysRange  []GenericKeyValue[Date, Stats] `json:"daysRange"`
	TotalRange Stats                          `json:"totalRange"`
}

type (

	// date underneath is a
	/* string of a time.Time format. "2006-01-02" */
	Date       string
	ScreenType string
	Category   string
)

type Stats struct {
	Active         float64        `json:"active"`
	Open           float64        `json:"open"`
	Inactive       float64        `json:"inactive"`
	ActiveTimeData []TimeInterval `json:"activeTimeData"`
}

type TimeInterval struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type ScreenTime struct {
	WindowID xproto.Window `json:"windowID"`
	AppName  string        `json:"appName"`
	Type     ScreenType    `json:"type"`
	Duration float64       `json:"duration"`
	Interval TimeInterval  `json:"interval"`
}

type Task struct {
	UUID       uuid.UUID  `json:"uuid"`
	AppName    string     `json:"appName"`
	AppCmdLine string     `json:"appCmdLine"`
	TaskTime   TaskTime   `json:"taskTime"`
	UI         UItextInfo `json:"ui"`
	Job        TaskType   `json:"job"`
}

type UItextInfo struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Notes    string `json:"notes"`
}

type TaskType string

type TaskTime struct {
	StartTime           time.Time `json:"startTime"`
	EndTime             time.Time `json:"endTime"`
	AlertTimesInMinutes [3]int    `json:"alertTimesInMinutes"`
	AlertSound          [3]bool   `json:"alertSound"`
}
