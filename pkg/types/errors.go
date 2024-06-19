package types

import "fmt"

var (
	ErrDeserialization = fmt.Errorf("error deserializing data")
	ErrSerialization   = fmt.Errorf("error serializing data")
	ErrLimitAppExist   = fmt.Errorf("limitApp task already exist")
)
