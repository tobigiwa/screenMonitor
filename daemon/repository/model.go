package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type appInfo struct {
	AppName           string
	Icon              []byte
	IsIconSet         bool
	Category          Category
	IsCategorySet     bool
	DesktopCategories []string
	ScreenStat        dailyAppScreenTime
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

type dailyScreentimeAnalytics map[date]dailyActiveScreentime

func (d *dailyScreentimeAnalytics) serialize() ([]byte, error) {
	var res bytes.Buffer
	encoded := gob.NewEncoder(&res)
	if err := encoded.Encode(d); err != nil {
		return nil, fmt.Errorf("%v:%w", err, ErrSerilization)
	}
	return res.Bytes(), nil
}

func (d *dailyScreentimeAnalytics) deserialize(data []byte) error {
	decoded := gob.NewDecoder(bytes.NewReader(data))
	if err := decoded.Decode(d); err != nil {
		return fmt.Errorf("%v:%w", err, ErrDeserilization)
	}
	return nil
}
