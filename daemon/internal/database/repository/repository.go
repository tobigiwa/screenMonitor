package repository

type IRepository interface {
	WriteUsage(data ScreenTime) error
	GetDay(Date) (DailyStat, error)

	GetWeek(string) (WeeklyStat, error)

	Close() error

	DeleteKey(key string) error
}
