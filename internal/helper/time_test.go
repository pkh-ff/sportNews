package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConverseToTimestampWithFormatRfc3339(t *testing.T) {
	timeStr := "2023-10-01T12:00:00Z"
	format := time.RFC3339
	expected := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)

	coverTime, err := ConverseToTimestamp(timeStr, format)

	require.NoError(t, err)
	require.Equal(t, expected, coverTime)
}

func TestConverseToTimestampWithDateOnly(t *testing.T) {
	timeStr := "2023-10-01"
	format := time.DateOnly
	expected := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)

	got, err := ConverseToTimestamp(timeStr, format)

	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestConverseToTimestampWithYYYYMMDD(t *testing.T) {
	timeStr := "20231001"
	format := "20060102"
	expected := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)

	got, err := ConverseToTimestamp(timeStr, format)

	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestConverseToTimestampWithFormatNoMatch(t *testing.T) {
	timeStr := "2023/10/01"
	format := "2006-01-02"

	got, err := ConverseToTimestamp(timeStr, format)

	require.Error(t, err)
	require.True(t, got.IsZero())
}

func TestConverseToTimestampWithInvalidString(t *testing.T) {
	timeStr := "invalid-time"
	format := time.RFC3339

	got, err := ConverseToTimestamp(timeStr, format)

	require.Error(t, err)
	assert.True(t, got.IsZero())
}

func TestConverseToTimestampWithEmptyInput(t *testing.T) {
	timeStr := ""
	format := time.RFC3339

	got, err := ConverseToTimestamp(timeStr, format)

	require.Error(t, err)
	assert.True(t, got.IsZero())
}

func TestConverseToTimestampWithInvalidFormat(t *testing.T) {
	timeStr := "2023/10/01"
	format := ""

	got, err := ConverseToTimestamp(timeStr, format)

	require.Error(t, err)
	assert.True(t, got.IsZero())
}
