package types

import "fmt"

var (
	NoMessage = Message{}
)

var (
	ErrDeserilization = fmt.Errorf("error deserializing data")
	ErrSerilization   = fmt.Errorf("error serializing data")
)
