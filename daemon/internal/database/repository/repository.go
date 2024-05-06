package repository

type IRepository interface {
	WriteUsage(data ScreenTime) error
	GetWeeklyScreenStats(ScreenType, string) ([]KeyValuePair, error)
	Close() error

	DeleteKey(key string) error
}
