package helper

import (
	"encoding/json"
	"fmt"
	"math"
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

func HrsAndMinute(hr float64) (int, int) {
	return int(hr), int(math.Round((hr - float64(int(hr))) * 60))
}

func UsageTimeInHrsMin(f float64) string {
	hrs, min := HrsAndMinute(f)
	return fmt.Sprintf("%dHrs:%dMin", hrs, min)
}
