package database

import "pkg/types"

type IRepository interface {
	WriteUsage(types.ScreenTime) error
	Close() error
	DeleteKey([]byte) error
	DeleteBucket(dbPrefix string) error
	GetDay(types.Date) (DailyStat, error)
	GetWeek(string) (WeeklyStat, error)
	AppWeeklyStat(appName string, anyDayInTheWeek types.Date) (types.AppRangeStat, error)
	AppMonthlyStat(appName, month, year string) (types.AppRangeStat, error)
	AppDateRangeStat(appName string, start, end types.Date) (types.AppRangeStat, error)
	GetAppIconCategoryAndCmdLine([]string) ([]types.AppIconCategoryAndCmdLine, error)
	UpdateOpertionOnBuCKET(dbPrefix string, opsFunc func([]byte) ([]byte, error)) error
}