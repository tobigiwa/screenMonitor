package service

import (
	"pkg/types"

	db "LiScreMon/cli/daemon/internal/database"

	"github.com/google/uuid"
)

type DatabaseInterface interface {
	WriteUsage(types.ScreenTime) error
	Close() error
	DeleteKey([]byte) error
	DeleteBucket(dbPrefix string) error
	GetDay(types.Date) (db.DailyStat, error)
	GetWeek(string) (db.WeeklyStat, error)
	AppWeeklyStat(appName string, anyDayInTheWeek types.Date) (types.AppRangeStat, error)
	AppMonthlyStat(appName, month, year string) (types.AppRangeStat, error)
	AppDateRangeStat(appName string, start, end types.Date) (types.AppRangeStat, error)
	GetAppIconCategoryAndCmdLine([]string) ([]types.AppIconCategoryAndCmdLine, error)
	UpdateOpertionOnBuCKET(dbPrefix string, opsFunc func([]byte) ([]byte, error)) error
	GetTaskByAppName(appName string) ([]types.Task, error)
	GetAllTask() ([]types.Task, error)
	RemoveTask(id uuid.UUID) error
	AddTask(task types.Task) error
	SetAppCategory(appName string, category types.Category) error
	GetAllACategories() ([]types.Category, error)
	UpdateAppInfoManually(key []byte, opsFunc func([]byte) ([]byte, error)) error
}
