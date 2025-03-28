package helper

import (
	"time"
)

// ConverseToTimestamp
// 字串時間轉換成 time.Time
func ConverseToTimestamp(timStr, format string) (time.Time, error) {
	t, err := time.Parse(format, timStr)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
