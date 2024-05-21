package types

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
	AppDetail       []ApplicationDetail `json:"appDetail"`
	IsError         bool                `json:"isError"`
	Error           error               `json:"error"`
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
	AppInfo AppIconAndCategory `json:"appInfo"`
	Usage   float64            `json:"usage"`
}
