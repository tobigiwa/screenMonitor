package store

import "fmt"

var (
	ErrAppKeyMismatch = fmt.Errorf("key error: app name mismatch")
	ErrDeserilization = fmt.Errorf("error deserializing data")
	ErrSerilization   = fmt.Errorf("error serializing data")
)
