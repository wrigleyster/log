package log

import (
	"testing"
	"time"
	"wlog/chrono"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	entry := NewLogEntry("init at 9:30").ParseTime()

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
		parseFrontDate().
		ParseTime()

	assert.Equal(t, "eod", entry.TaskName)
	assert.Equal(t, "Friday", entry.Time.Weekday().String())
	assert.Equal(t, friday.Hour(), entry.Time.Hour())
	assert.Equal(t, friday.Minute(), entry.Time.Minute())
	assert.Equal(t, friday.Year(), entry.Time.Year())
	assert.Equal(t, friday.Month(), entry.Time.Month())
	assert.Equal(t, friday.Day(), entry.Time.Day())
}
func TestParseFrontDateTime(t *testing.T) {
	friday := time.Date(2023, 6, 2, 16, 45, 0, 0, time.Local)
	today := time.Date(2023, 6, 5, 8, 17, 0, 0, time.Local)
	entry := Entry{today, "friday 16:45 eod", ""}.
		parseDate().
		parseFrontDate().
		ParseTime()

	assert.Equal(t, "eod", entry.TaskName)
	assert.Equal(t, "Friday", entry.Time.Weekday().String())
	assert.Equal(t, friday.Hour(), entry.Time.Hour())
	assert.Equal(t, friday.Minute(), entry.Time.Minute())
	assert.Equal(t, friday.Year(), entry.Time.Year())
	assert.Equal(t, friday.Month(), entry.Time.Month())
	assert.Equal(t, friday.Day(), entry.Time.Day())
}
func TestParseFullFrontDateTime(t *testing.T) {
	friday := time.Date(2023, 6, 2, 16, 45, 0, 0, time.Local)
	today := time.Date(2023, 6, 5, 8, 17, 0, 0, time.Local)
	entry := Entry{today, "2023.6.2 16:45 eod", ""}.
		parseDate().
		parseFrontDate().
		ParseTime()

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
	r := chrono.RelativeDate(today, "friday")
	assert.Equal(t, "Friday", r.Weekday().String())
	assert.Equal(t, friday.Year(), r.Year())
	assert.Equal(t, friday.Month(), r.Month())
	assert.Equal(t, friday.Day(), r.Day())
}

func TestParseDate_yesterday(t *testing.T) {
	entry := NewLogEntry("working on cool stuff at 8:40 yesterday").
		parseDate().
		parseFrontDate().
		ParseTime()

	yesterday := time.Now().AddDate(0, 0, -1)
	assert.Equal(t, "working on cool stuff", entry.TaskName)
	assert.Equal(t, 8, entry.Time.Hour())
	assert.Equal(t, 40, entry.Time.Minute())
	assert.Equal(t, yesterday.Year(), entry.Time.Year())
	assert.Equal(t, yesterday.Month(), entry.Time.Month())
	assert.Equal(t, yesterday.Day(), entry.Time.Day())
}

func TestParseDate_y(t *testing.T) {
	entry := NewLogEntry("working on cool stuff at 8:40 y").
		parseDate().
		parseFrontDate().
		ParseTime()

	yesterday := time.Now().AddDate(0, 0, -1)
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
		parseFrontDate().
		ParseTime()

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

	monday := chrono.RelativeDate(today, "monday")
	assert.Equal(t, "Monday", monday.Weekday().String())
	assert.True(t, today.After(monday))

	tuesday := chrono.RelativeDate(today, "tuesday")
	assert.Equal(t, "Tuesday", tuesday.Weekday().String())
	assert.True(t, today.After(tuesday))

	wednesday := chrono.RelativeDate(today, "wednesday")
	assert.Equal(t, "Wednesday", wednesday.Weekday().String())
	assert.True(t, today.After(wednesday))

	thursday := chrono.RelativeDate(today, "thursday")
	assert.Equal(t, "Thursday", thursday.Weekday().String())
	assert.True(t, today.After(thursday))

	friday := chrono.RelativeDate(today, "friday")
	assert.Equal(t, "Friday", friday.Weekday().String())
	assert.True(t, today.After(friday))

	saturday := chrono.RelativeDate(today, "saturday")
	assert.Equal(t, "Saturday", saturday.Weekday().String())
	assert.True(t, today.After(saturday))

	sunday := chrono.RelativeDate(today, "sunday")
	assert.Equal(t, "Sunday", sunday.Weekday().String())
	assert.Equal(t, today, sunday)
}
func TestRelativeDate2(t *testing.T) {
	today := time.Date(2023, 6, 5, 8, 17, 0, 0, time.Local)
	assert.Equal(t, "Sunday", time.Weekday(0).String())
	assert.Equal(t, time.Weekday(1), today.Weekday())

	monday := chrono.RelativeDate(today, "monday")
	assert.Equal(t, "Monday", monday.Weekday().String())
	assert.Equal(t, today, monday, "today is monday")

	tuesday := chrono.RelativeDate(today, "tuesday")
	assert.Equal(t, "Tuesday", tuesday.Weekday().String())
	assert.True(t, today.After(tuesday), "tuesday is later")

	wednesday := chrono.RelativeDate(today, "wednesday")
	assert.Equal(t, "Wednesday", wednesday.Weekday().String())
	assert.True(t, today.After(wednesday), "wednesday is later")

	thursday := chrono.RelativeDate(today, "thursday")
	assert.Equal(t, "Thursday", thursday.Weekday().String())
	assert.True(t, today.After(thursday), "thursday is later")

	friday := chrono.RelativeDate(today, "friday")
	assert.Equal(t, "Friday", friday.Weekday().String())
	assert.True(t, today.After(friday), "friday is later")

	saturday := chrono.RelativeDate(today, "saturday")
	assert.Equal(t, "Saturday", saturday.Weekday().String())
	assert.True(t, today.After(saturday), "saturday is later")

	sunday := chrono.RelativeDate(today, "sunday")
	assert.Equal(t, "Sunday", sunday.Weekday().String())
	assert.True(t, today.After(sunday))
}

func TestParseTaskId(t *testing.T) {
	entry := Entry{time.Now(), "I want to work on SFFEAT001234", ""}

	entry = entry.parseExtId()

	assert.Equal(t, "SFFEAT001234", entry.ExtId)
}

func TestRelativeDateMonday(t *testing.T) {
	monday := time.Date(2023, 5, 22, 10, 54, 0, 0, time.Local)
	tuesday := time.Date(2023, 5, 23, 10, 54, 0, 0, time.Local)
	wednesday := time.Date(2023, 5, 24, 10, 54, 0, 0, time.Local)
	thursday := time.Date(2023, 5, 25, 10, 54, 0, 0, time.Local)
	friday := time.Date(2023, 5, 26, 10, 54, 0, 0, time.Local)

	assert.Equal(t, "Monday", monday.Weekday().String())
	assert.Equal(t, chrono.Day(monday), chrono.Day(chrono.RelativeDate(tuesday, "monday")))
	assert.Equal(t, chrono.Day(monday), chrono.Day(chrono.RelativeDate(wednesday, "monday")))
	assert.Equal(t, chrono.Day(monday), chrono.Day(chrono.RelativeDate(thursday, "monday")))
	assert.Equal(t, chrono.Day(monday), chrono.Day(chrono.RelativeDate(friday, "monday")))

	assert.Equal(t, "Tuesday", tuesday.Weekday().String())
	assert.Equal(t, chrono.Day(tuesday), chrono.Day(chrono.RelativeDate(tuesday, "tuesday")))
	assert.Equal(t, chrono.Day(tuesday), chrono.Day(chrono.RelativeDate(wednesday, "tuesday")))
	assert.Equal(t, chrono.Day(tuesday), chrono.Day(chrono.RelativeDate(thursday, "tuesday")))
	assert.Equal(t, chrono.Day(tuesday), chrono.Day(chrono.RelativeDate(friday, "tuesday")))
}
func TestAbsoluteDate(t *testing.T) {
	monday := time.Date(2023, 5, 22, 10, 54, 0, 0, time.Local)
	assert.Equal(t, "Monday", monday.Weekday().String())
	assert.Equal(t, chrono.Day(monday), chrono.Day(chrono.AbsoluteDate(monday, "2023.5.22")))
}
func TestFrontDate(t *testing.T) {
	monday := time.Date(2023, 5, 22, 10, 54, 0, 0, time.Local)
	entry := Entry{time.Now(), "2023.5.22 10:54 ducks are winners", ""}

	entry = entry.parseFrontDate()
	
	assert.Equal(t, monday.Day(), entry.Time.Day())
	assert.Equal(t, monday.Month(), entry.Time.Month())
	assert.Equal(t, monday.Year(), entry.Time.Year())
	assert.Equal(t, monday.YearDay(), entry.Time.YearDay())
	assert.Equal(t, "10:54 ducks are winners", entry.TaskName)
}
func TestFrontTime(t *testing.T) {
	monday := time.Date(2023, 5, 22, 10, 54, 0, 0, time.Local)
	entry := Entry{chrono.Date(monday).At(0,0), "10:54 ducks are winners", ""}
	

	entry = entry.ParseTime()
	
	assert.Equal(t, monday.Hour(), entry.Time.Hour())
	assert.Equal(t, monday.Minute(), entry.Time.Minute())
	assert.Equal(t, monday.Location(), entry.Time.Location())
	assert.Equal(t, "ducks are winners", entry.TaskName)
	assert.Equal(t, monday, entry.Time)
}
