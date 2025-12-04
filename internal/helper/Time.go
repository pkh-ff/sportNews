package helper

import (
	"sportNews/pkg/log"
	"time"

	"go.uber.org/zap"
)

// ConverseToTimestamp
// 字串時間轉換成 time.Time
func ConverseToTimestamp(timeStr, format string) (time.Time, error) {
	log.Info("ConverseToTimestamp: Converting time string to timestamp",
		zap.String("timeStr", timeStr),
		zap.String("format", format),
	)

	t, err := time.Parse(format, timeStr)
	if err != nil {
		log.Error("ConverseToTimestamp: Failed to parse time string",
			zap.String("timeStr", timeStr),
			zap.String("format", format),
			zap.Error(err))
		return time.Time{}, err
	}

	log.Info("ConverseToTimestamp: Successfully converted time string",
		zap.String("timeStr", timeStr),
		zap.Time("parsedTime", t),
	)
	return t, nil
}
