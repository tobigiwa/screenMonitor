package store

type Repository interface {
	WriteUsage(data ScreenTime) error
}
