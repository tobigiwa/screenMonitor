package helper

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"pkg/types"
)

func Encode[T any](tyPe T) ([]byte, error) {
	var r bytes.Buffer
	encoded := gob.NewEncoder(&r)
	if err := encoded.Encode(tyPe); err != nil {
		return nil, fmt.Errorf("%v:%w", err, types.ErrSerilization)
	}
	return r.Bytes(), nil
}

func Decode[T any](data []byte) (T, error) {
	var t, result T
	decoded := gob.NewDecoder(bytes.NewReader(data))
	if err := decoded.Decode(&result); err != nil {
		return t, fmt.Errorf("%v:%w", err, types.ErrDeserilization)
	}
	return result, nil
}
