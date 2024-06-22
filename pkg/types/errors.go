package types

import "fmt"

var (
	ErrDeserialization = fmt.Errorf("error deserializing data")
	ErrSerialization   = fmt.Errorf("error serializing data")
	ErrLimitAppExist   = fmt.Errorf("limitApp task already exist")

	ErrDeletingTask         = fmt.Errorf("err deleting old task")
	ErrTaskMangerNotStarted = fmt.Errorf("taskManager could not be started")

	ErrZeroValueTask = fmt.Errorf("sent task cannot be empty struct")
)
