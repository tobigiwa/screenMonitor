package helper

import (
	"encoding/json"
	"fmt"
	"pkg/types"
)

func Encode[T any](tyPe T) ([]byte, error) {
	encoded, err := json.Marshal(tyPe)
	if err != nil {
		return nil, fmt.Errorf("%v:%w", err, types.ErrSerialization)
	}
	return encoded, nil
}

func Decode[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, fmt.Errorf("%v:%w", err, types.ErrDeserialization)
	}
	return result, nil
}
