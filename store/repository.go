package store

type IRepository interface {
	WriteUsage(data ScreenTime) error
}
