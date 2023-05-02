package formatter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wlog/log"
)

func TestFormat(t *testing.T) {
	date := time.Date(2023, 4, 23, 14, 37, 0, 0, time.Local)
	entries := []log.Entry{
		{log.Date(date).At(14, 30), "Testing the timelogger", "SFFEAT0000003"},
		{log.Date(date).At(12, 00), "Implementing the timelogger", "SFFEAT0000002"},
		{log.Date(date).At(9, 30), "Designing the timelogger", "SFFEAT0000001"},
	}
	actual := Format(entries)

	expected := "Sunday 23 April\n" +
		" 14:30 SFFEAT0000003 Testing the timelogger\n" +
		" 12:00 SFFEAT0000002 Implementing the timelogger\n" +
		" 09:30 SFFEAT0000001 Designing the timelogger"

	assert.Equal(t, expected, actual)
}

func TestFormatDurations_withEOD(t *testing.T) {
	date := time.Date(2023, 4, 23, 14, 37, 0, 0, time.Local)
	entries := []log.Entry{
		{log.Date(date).At(16, 30), "EOD", ""},
		{log.Date(date).At(14, 30), "Testing the timelogger", "SFFEAT0000003"},
		{log.Date(date).At(12, 00), "Implementing the timelogger", "SFFEAT0000002"},
		{log.Date(date).At(11, 31), "lunch", ""},
		{log.Date(date).At(9, 30), "Designing the timelogger", "SFFEAT0000001"},
	}
	actual := FormatDurations(entries, time.Now())

	expected := "Sunday 23 April\n" +
		" 2h 00m SFFEAT0000003 Testing the timelogger\n" +
		" 2h 30m SFFEAT0000002 Implementing the timelogger\n" +
		" 0h 29m lunch\n" +
		" 2h 01m SFFEAT0000001 Designing the timelogger"

	assert.Equal(t, expected, actual)
}

func TestFormatDurations_withoutEOD(t *testing.T) {
	date := time.Date(2023, 4, 23, 16, 37, 0, 0, time.Local)
	entries := []log.Entry{
		{log.Date(date).At(14, 30), "Testing the timelogger", "SFFEAT0000003"},
		{log.Date(date).At(12, 00), "Implementing the timelogger", "SFFEAT0000002"},
		{log.Date(date).At(11, 31), "lunch", ""},
		{log.Date(date).At(9, 30), "Designing the timelogger", "SFFEAT0000001"},
	}
	actual := FormatDurations(entries, date)

	expected := "Sunday 23 April\n" +
		"+2h 07m SFFEAT0000003 Testing the timelogger\n" +
		" 2h 30m SFFEAT0000002 Implementing the timelogger\n" +
		" 0h 29m lunch\n" +
		" 2h 01m SFFEAT0000001 Designing the timelogger"

	assert.Equal(t, expected, actual)
}
