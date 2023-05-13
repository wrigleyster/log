package manipulation

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wlog/log"
)

func TestAccumulateWithoutEOD(t *testing.T) {
	date := log.Date(time.Now()).At(15, 00)
	entries := []log.Entry{
		{log.Date(date).At(9, 30), "Designing the timelogger", "SFFEAT0000001"},
		{log.Date(date).At(12, 00), "Implementing the timelogger", "SFFEAT0000002"},
		{log.Date(date).At(14, 30), "Testing the timelogger", "SFFEAT0000003"},
	}
	actual := Accumulate(entries, date)

	expected := Total([]DayTotal{
		{DurationOf(5, 30), GetDay(date), []TaskTotal{
			{DurationOf(2, 30), "Designing the timelogger", "SFFEAT0000001", false},
			{DurationOf(2, 30), "Implementing the timelogger", "SFFEAT0000002", false},
			{DurationOf(0, 30), "Testing the timelogger", "SFFEAT0000003", true},
		}},
	})

	assert.Equal(t, expected, actual)
}
func TestAccumulateWithEOD(t *testing.T) {
	date := log.Date(time.Now()).At(15, 00)
	entries := []log.Entry{
		{log.Date(date).At(9, 30), "Designing the timelogger", "SFFEAT0000001"},
		{log.Date(date).At(12, 00), "Implementing the timelogger", "SFFEAT0000002"},
		{log.Date(date).At(14, 30), "Testing the timelogger", "SFFEAT0000003"},
		{log.Date(date).At(15, 00), "eod", ""},
	}
	actual := Accumulate(entries, date)

	expected := Total([]DayTotal{
		{DurationOf(5, 30), GetDay(date), []TaskTotal{
			{DurationOf(2, 30), "Designing the timelogger", "SFFEAT0000001", false},
			{DurationOf(2, 30), "Implementing the timelogger", "SFFEAT0000002", false},
			{DurationOf(0, 30), "Testing the timelogger", "SFFEAT0000003", false},
			{DurationOf(0, 00), "eod", "", false},
		}},
	})

	assert.Equal(t, expected, actual)
}
func TestAccumulateWithEodAndNewDay(t *testing.T) {
	now := log.Date(time.Now()).At(15, 00)
	yesterday := now.Add(-24 * time.Hour)
	entries := []log.Entry{
		{log.Date(yesterday).At(9, 30), "Designing the timelogger", "SFFEAT0000001"},
		{log.Date(yesterday).At(12, 00), "Implementing the timelogger", "SFFEAT0000002"},
		{log.Date(yesterday).At(14, 30), "Testing the timelogger", "SFFEAT0000003"},
		{log.Date(yesterday).At(15, 00), "eod", ""},
		{log.Date(now).At(01, 00), "early start", ""},
	}
	actual := Accumulate(entries, now)

	expected := Total([]DayTotal{
		{DurationOf(5, 30), GetDay(yesterday), []TaskTotal{
			{DurationOf(2, 30), "Designing the timelogger", "SFFEAT0000001", false},
			{DurationOf(2, 30), "Implementing the timelogger", "SFFEAT0000002", false},
			{DurationOf(0, 30), "Testing the timelogger", "SFFEAT0000003", false},
			{DurationOf(0, 00), "eod", "", false},
		}},
		{DurationOf(14, 00), GetDay(now), []TaskTotal{
			{DurationOf(14, 00), "early start", "", true},
		}},
	})

	assert.Equal(t, expected, actual)
}
