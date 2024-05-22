package database

import "pkg/types"

type IRepository interface {
	WriteUsage(data ScreenTime) error
	Close() error
	DeleteKey(key string) error
	GetDay(Date) (DailyStat, error)
	GetWeek(string) (WeeklyStat, error)
	GetAppIconAndCategory(appNames []string) ([]types.AppIconAndCategory, error)
}
