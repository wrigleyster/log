package log

import (
	"github.com/stretchr/testify/assert"
	"wlog/chrono"
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

func TestParseDateTime(t *testing.T) {
	friday := time.Date(2023, 6, 2, 16, 45, 0, 0, time.Local)
	today := time.Date(2023, 6, 5, 8, 17, 0, 0, time.Local)
	entry := Entry{today, "eod at 16:45 friday", ""}.
		parseDate().
		parseTime()

	assert.Equal(t, "eod", entry.TaskName)
	assert.Equal(t, "Friday", entry.Time.Weekday().String())
	assert.Equal(t, friday.Hour(), entry.Time.Hour())
	assert.Equal(t, friday.Minute(), entry.Time.Minute())
	assert.Equal(t, friday.Year(), entry.Time.Year())
	assert.Equal(t, friday.Month(), entry.Time.Month())
	assert.Equal(t, friday.Day(), entry.Time.Day())
}
func TestRelativeDateTime(t *testing.T) {
	friday := time.Date(2023, 6, 2, 16, 45, 0, 0, time.Local)
	today := time.Date(2023, 6, 5, 8, 17, 0, 0, time.Local)
	r := relativeDate(today, "friday")
	assert.Equal(t, "Friday", r.Weekday().String())
	assert.Equal(t, friday.Year(), r.Year())
	assert.Equal(t, friday.Month(), r.Month())
	assert.Equal(t, friday.Day(), r.Day())
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
	assert.True(t, monday.Before(entry.Time))
}

func TestRelativeDate(t *testing.T) {
	today := time.Date(2023, 4, 30, 11, 59, 0, 0, time.Local)
	assert.Equal(t, "Sunday", time.Weekday(0).String())
	assert.Equal(t, time.Weekday(0), today.Weekday())

	monday := relativeDate(today, "monday")
	assert.Equal(t, "Monday", monday.Weekday().String())
	assert.True(t, today.After(monday))

	tuesday := relativeDate(today, "tuesday")
	assert.Equal(t, "Tuesday", tuesday.Weekday().String())
	assert.True(t, today.After(tuesday))

	wednesday := relativeDate(today, "wednesday")
	assert.Equal(t, "Wednesday", wednesday.Weekday().String())
	assert.True(t, today.After(wednesday))

	thursday := relativeDate(today, "thursday")
	assert.Equal(t, "Thursday", thursday.Weekday().String())
	assert.True(t, today.After(thursday))

	friday := relativeDate(today, "friday")
	assert.Equal(t, "Friday", friday.Weekday().String())
	assert.True(t, today.After(friday))

	saturday := relativeDate(today, "saturday")
	assert.Equal(t, "Saturday", saturday.Weekday().String())
	assert.True(t, today.After(saturday))

	sunday := relativeDate(today, "sunday")
	assert.Equal(t, "Sunday", sunday.Weekday().String())
	assert.Equal(t, today, sunday)
}
func TestRelativeDate2(t *testing.T) {
	today := time.Date(2023, 6, 5, 8, 17, 0, 0, time.Local)
	assert.Equal(t, "Sunday", time.Weekday(0).String())
	assert.Equal(t, time.Weekday(1), today.Weekday())

	monday := relativeDate(today, "monday")
	assert.Equal(t, "Monday", monday.Weekday().String())
	assert.Equal(t, today, monday, "today is monday")

	tuesday := relativeDate(today, "tuesday")
	assert.Equal(t, "Tuesday", tuesday.Weekday().String())
	assert.True(t, today.After(tuesday), "tuesday is later")

	wednesday := relativeDate(today, "wednesday")
	assert.Equal(t, "Wednesday", wednesday.Weekday().String())
	assert.True(t, today.After(wednesday), "wednesday is later")

	thursday := relativeDate(today, "thursday")
	assert.Equal(t, "Thursday", thursday.Weekday().String())
	assert.True(t, today.After(thursday), "thursday is later")

	friday := relativeDate(today, "friday")
	assert.Equal(t, "Friday", friday.Weekday().String())
	assert.True(t, today.After(friday), "friday is later")

	saturday := relativeDate(today, "saturday")
	assert.Equal(t, "Saturday", saturday.Weekday().String())
	assert.True(t, today.After(saturday), "saturday is later")

	sunday := relativeDate(today, "sunday")
	assert.Equal(t, "Sunday", sunday.Weekday().String())
	assert.True(t, today.After(sunday))
}

func TestParseTaskId(t *testing.T) {
	entry := Entry{time.Now(), "I want to work on SFFEAT001234", ""}

	entry = entry.parseTaskId()

	assert.Equal(t, "SFFEAT001234", entry.TaskId)
}

func TestRelativeDateMonday(t *testing.T) {
	monday := time.Date(2023, 5, 22, 10, 54, 0, 0, time.Local)
	tuesday := time.Date(2023, 5, 23, 10, 54, 0, 0, time.Local)
	wednesday := time.Date(2023, 5, 24, 10, 54, 0, 0, time.Local)
	thursday := time.Date(2023, 5, 25, 10, 54, 0, 0, time.Local)
	friday := time.Date(2023, 5, 26, 10, 54, 0, 0, time.Local)

	assert.Equal(t, "Monday", monday.Weekday().String())
	assert.Equal(t, chrono.Day(monday), chrono.Day(relativeDate(tuesday, "monday")))
	assert.Equal(t, chrono.Day(monday), chrono.Day(relativeDate(wednesday, "monday")))
	assert.Equal(t, chrono.Day(monday), chrono.Day(relativeDate(thursday, "monday")))
	assert.Equal(t, chrono.Day(monday), chrono.Day(relativeDate(friday, "monday")))

	assert.Equal(t, "Tuesday", tuesday.Weekday().String())
	assert.Equal(t, chrono.Day(tuesday), chrono.Day(relativeDate(tuesday, "tuesday")))
	assert.Equal(t, chrono.Day(tuesday), chrono.Day(relativeDate(wednesday, "tuesday")))
	assert.Equal(t, chrono.Day(tuesday), chrono.Day(relativeDate(thursday, "tuesday")))
	assert.Equal(t, chrono.Day(tuesday), chrono.Day(relativeDate(friday, "tuesday")))
}
