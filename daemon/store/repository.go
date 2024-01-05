package store

type IRepository interface {
	WriteUsage(data ScreenTime) error
	GetWeeklyScreenStats(s ScreenType) (map[date]float64, error)
	Close() error

	DeleteKey(key string) error
}
