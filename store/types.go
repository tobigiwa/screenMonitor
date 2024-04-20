package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type ScreenType string

const (
	Active  ScreenType = "active"
	Passive ScreenType = "passive"
)

// ScreenTime represents the time spent on a particular app.
type ScreenTime struct {
	AppName  string
	Type     ScreenType
	Duration float64
	Interval TimeInterval
}

type date string

type dailyAppScreenTime map[date]stats

type TimeInterval struct {
	Start time.Time
	End   time.Time
}

type stats struct {
	Active         float64
	Open           float64
	Inactive       float64
	ActiveTimeData []TimeInterval
}

type appInfo struct {
	AppName    string
	Icon       []byte
	Category   string
	ScreenStat dailyAppScreenTime
}

func (ap *appInfo) serialize() ([]byte, error) {
	var res bytes.Buffer
	encoded := gob.NewEncoder(&res)
	if err := encoded.Encode(ap); err != nil {
		return nil, fmt.Errorf("%v:%w", err, ErrSerilization)
	}
	return res.Bytes(), nil
}

func (ap *appInfo) deserialize(data []byte) error {
	decoded := gob.NewDecoder(bytes.NewReader(data))
	if err := decoded.Decode(ap); err != nil {
		return fmt.Errorf("%v:%w", err, ErrDeserilization)
	}
	return nil
}
