package types

import "fmt"

var (
	NoMessage = Message{}
)

var (
	ErrDeserialization = fmt.Errorf("error deserializing data")
	ErrSerialization   = fmt.Errorf("error serializing data")
)
