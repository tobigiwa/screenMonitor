package repository

type IRepository interface {
	WriteUsage(data ScreenTime) error
	GetWeeklyScreenStats(ScreenType, string) ([7]KeyValuePair, error)
	GetDay(Date) (DailyStat, error)

	GetWeek(anyDayInTheWeek Date) (WeeklyStat, error)

	Close() error

	DeleteKey(key string) error
}
