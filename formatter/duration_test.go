package formatter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wlog/chrono"
	"wlog/log"
	"wlog/manipulation"
)

func TestFormatDurations_withEOD(t *testing.T) {
	date := time.Date(2023, 4, 23, 14, 37, 0, 0, time.Local)
	entries := []log.Entry{
		{Time: chrono.Date(date).At(9, 30), TaskName: "Designing the timelogger", TaskId: "SFFEAT0000001"},
		{Time: chrono.Date(date).At(11, 31),TaskName: "lunch", TaskId: ""},
		{Time: chrono.Date(date).At(12, 00),TaskName: "Implementing the timelogger", TaskId: "SFFEAT0000002"},
		{Time: chrono.Date(date).At(14, 30),TaskName: "Testing the timelogger", TaskId: "SFFEAT0000003"},
		{Time: chrono.Date(date).At(16, 30),TaskName: "EOD", TaskId: ""},
	}
	actual := DurationView(manipulation.Aggregate(entries, time.Now())).Format(Ascending)

	expected := "Sunday 23 April, total: 7h 00m\n" +
		" 2h 01m SFFEAT0000001 Designing the timelogger\n" +
		" 0h 29m lunch\n" +
		" 2h 30m SFFEAT0000002 Implementing the timelogger\n" +
		" 2h 00m SFFEAT0000003 Testing the timelogger"

	assert.Equal(t, expected, actual)
}

func TestFormatDurations_withoutEOD(t *testing.T) {
	date := time.Date(2023, 4, 23, 16, 37, 0, 0, time.Local)
	entries := []log.Entry{
		{Time: chrono.Date(date).At(9, 30), TaskName: "Designing the timelogger", TaskId: "SFFEAT0000001"},
		{Time: chrono.Date(date).At(11, 31),TaskName: "lunch", TaskId: ""},
		{Time: chrono.Date(date).At(12, 00),TaskName: "Implementing the timelogger", TaskId: "SFFEAT0000002"},
		{Time: chrono.Date(date).At(14, 30),TaskName: "Testing the timelogger", TaskId: "SFFEAT0000003"},
	}
	actual := DurationView(manipulation.Aggregate(entries, date)).Format(Ascending)

	expected := "Sunday 23 April, total: 7h 07m\n" +
		" 2h 01m SFFEAT0000001 Designing the timelogger\n" +
		" 0h 29m lunch\n" +
		" 2h 30m SFFEAT0000002 Implementing the timelogger\n" +
		"+2h 07m SFFEAT0000003 Testing the timelogger"

	assert.Equal(t, expected, actual)
}
func TestFormatTotal(t *testing.T) {
	date := time.Date(2023, 4, 23, 14, 37, 0, 0, time.Local)
	earlier := time.Date(2023, 4, 22, 14, 37, 0, 0, time.Local)
	entries := []log.Entry{
		{Time: chrono.Date(earlier).At(9, 00), TaskName: "working", TaskId: "SFFEAT0000001"},
		{Time: chrono.Date(earlier).At(17, 00),TaskName: "eod", TaskId: ""},
		{Time: chrono.Date(date).At(9, 30),  TaskName: "Designing the timelogger", TaskId: "SFFEAT0000001"},
		{Time: chrono.Date(date).At(12, 00), TaskName: "Implementing the timelogger", TaskId: "SFFEAT0000002"},
		{Time: chrono.Date(date).At(14, 30), TaskName: "Testing the timelogger", TaskId: "SFFEAT0000003"},
	}
	actual := DurationView(manipulation.Aggregate(entries, date)).Format(Ascending)

	expected := "Saturday 22 April, total: 8h 00m\n" +
		" 8h 00m SFFEAT0000001 working\n" +
		"Sunday 23 April, total: 5h 07m\n" +
		" 2h 30m SFFEAT0000001 Designing the timelogger\n" +
		" 2h 30m SFFEAT0000002 Implementing the timelogger\n" +
		"+0h 07m SFFEAT0000003 Testing the timelogger"

	assert.Equal(t, expected, actual)
}
func TestFormatTotal2(t *testing.T) {
	date := time.Date(2023, 5, 13, 1, 02, 0, 0, time.Local)
	earlier := time.Date(2023, 5, 12, 14, 02, 0, 0, time.Local)
	entries := []log.Entry{
		{Time: chrono.Date(earlier).At(9, 30),  TaskName: "dsu", TaskId: ""},
		{Time: chrono.Date(earlier).At(10, 00), TaskName: "horse riding", TaskId: ""},
		{Time: chrono.Date(earlier).At(12, 00), TaskName: "eod", TaskId: ""},
		{Time: chrono.Date(date).At(0, 43), TaskName: "early start", TaskId: ""},
	}
	actual := DurationView(manipulation.Aggregate(entries, date)).Format(Descending)

	expected := "Saturday 13 May, total: 0h 19m\n" +
		"+0h 19m early start\n" +
		"Friday 12 May, total: 2h 30m\n" +
		" 2h 00m horse riding\n" +
		" 0h 30m dsu"

	assert.Equal(t, expected, actual)
}
func TestFormatTotal3(t *testing.T) {
	date := time.Date(2023, 5, 13, 10, 02, 0, 0, time.Local)
	earlier := time.Date(2023, 5, 12, 14, 02, 0, 0, time.Local)
	entries := []log.Entry{
		{Time: chrono.Date(earlier).At(9, 30),  TaskName: "dsu", TaskId: ""},
		{Time: chrono.Date(earlier).At(9, 45),  TaskName: "check email", TaskId: ""},
		{Time: chrono.Date(earlier).At(10, 00), TaskName: "horse riding", TaskId: ""},
		{Time: chrono.Date(earlier).At(11, 45), TaskName: "check email", TaskId: ""},
		{Time: chrono.Date(earlier).At(12, 00), TaskName: "eod", TaskId: ""},
		{Time: chrono.Date(date).At(9, 30), TaskName: "check email", TaskId: ""},
	}
	actual := DurationView(manipulation.Aggregate(entries, date)).Format(Descending)

	expected := "Saturday 13 May, total: 0h 32m\n" +
		"+0h 32m check email\n" +
		"Friday 12 May, total: 2h 30m\n" +
		" 0h 30m check email\n" +
		" 1h 45m horse riding\n" +
		" 0h 15m dsu"

	assert.Equal(t, expected, actual)
}
