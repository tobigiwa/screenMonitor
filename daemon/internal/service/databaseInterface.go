package service

import (
	db "LiScreMon/daemon/internal/database"

	"github.com/google/uuid"
	"utils"
)

type DatabaseInterface interface {
	WriteUsage(utils.ScreenTime) error
	DeleteKey([]byte) error
	DeleteBucket(dbPrefix string) error
	GetDay(utils.Date) (db.DailyStat, error)
	GetWeek(utils.Date) (db.WeeklyStat, error)
	AppWeeklyStat(appName string, anyDayInTheWeek utils.Date) (utils.AppRangeStat, error)
	AppMonthlyStat(appName, month, year string) (utils.AppRangeStat, error)
	AppDateRangeStat(appName string, start, end utils.Date) (utils.AppRangeStat, error)
	GetAppIconCategoryAndCmdLine([]string) ([]utils.AppIconCategoryAndCmdLine, error)
	UpdateOpertionOnBuCKET(dbPrefix string, opsFunc func([]byte) ([]byte, error)) error
	GetTaskByAppName(appName string) ([]utils.Task, error)
	GetAllTask() ([]utils.Task, error)
	RemoveTask(id uuid.UUID) error
	AddTask(task utils.Task) error
	SetAppCategory(appName string, category utils.Category) error
	GetAllACategories() ([]utils.Category, error)
	UpdateAppInfoManually(key []byte, opsFunc func([]byte) ([]byte, error)) error
	GetAllApp() ([]utils.AppIconCategoryAndCmdLine, error)
	GetTaskByUUID(taskID uuid.UUID) (utils.Task, error)
}
