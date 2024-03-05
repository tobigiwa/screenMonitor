package store

type Repository interface {
	WriteUsuage(data ScreenTime) error
}
