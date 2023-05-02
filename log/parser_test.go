package log

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	entry := NewLogEntry("init at 9:30").parseTime()

	assert.Equal(t, "init", entry.TaskName)
	assert.Equal(t, 9, entry.Time.Hour())
	assert.Equal(t, 30, entry.Time.Minute())
	assert.Equal(t, time.Now().Year(), entry.Time.Year())
	assert.Equal(t, time.Now().Month(), entry.Time.Month())
	assert.Equal(t, time.Now().Day(), entry.Time.Day())
}

func TestParseDate_yesterday(t *testing.T) {
	entry := NewLogEntry("working on cool stuff at 8:40 yesterday").
		parseDate().
		parseTime()

	yesterday := time.Now().Add(-time.Hour * 24)
	assert.Equal(t, "working on cool stuff", entry.TaskName)
	assert.Equal(t, 8, entry.Time.Hour())
	assert.Equal(t, 40, entry.Time.Minute())
	assert.Equal(t, yesterday.Year(), entry.Time.Year())
	assert.Equal(t, yesterday.Month(), entry.Time.Month())
	assert.Equal(t, yesterday.Day(), entry.Time.Day())
}

func TestParseDate_monday(t *testing.T) {
	entry := NewLogEntry("working on cool stuff at 8:40 monday").
		parseDate().
		parseTime()

	monday := time.Now().Truncate(time.Hour * 24 * 7)
	assert.Equal(t, "working on cool stuff", entry.TaskName)
	assert.Equal(t, 8, entry.Time.Hour())
	assert.Equal(t, 40, entry.Time.Minute())
	assert.Equal(t, monday.Year(), entry.Time.Year())
	assert.Equal(t, monday.Month(), entry.Time.Month())
	assert.Equal(t, monday.Day(), entry.Time.Day())
}

func TestRelativeDate(t *testing.T) {
	today := time.Date(2023, 4, 30, 11, 59, 0, 0, time.Local)
	assert.Equal(t, "Sunday", time.Weekday(0).String())
	assert.Equal(t, time.Weekday(0), today.Weekday())

	tuesday := relativeDate(today, "tuesday")
	assert.Equal(t, "Tuesday", tuesday.Weekday().String())
	assert.True(t, today.After(tuesday))

	monday := relativeDate(today, "monday")
	assert.Equal(t, "Monday", monday.Weekday().String())
	assert.True(t, today.After(monday))

	sunday := relativeDate(today, "sunday")
	assert.Equal(t, "Sunday", sunday.Weekday().String())
	assert.Equal(t, today, sunday)
}

func TestParseTaskId(t *testing.T) {
	entry := Entry{time.Now(), "I want to work on SFFEAT001234", ""}

	entry = entry.parseTaskId()

	assert.Equal(t, "SFFEAT001234", entry.TaskId)
}
