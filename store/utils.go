package store

import (
	"fmt"
	"time"
)

func Key() string {
	now := time.Now()
	day, month, year := now.Day(), now.Month(), time.Now().Year()
	return fmt.Sprintf("%d:%d:%d", day, month, year)
}
