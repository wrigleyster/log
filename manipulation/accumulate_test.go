package manipulation

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wlog/chrono"
	"wlog/log"
)

func TestAccumulateWithoutEOD(t *testing.T) {
	date := chrono.Date(time.Now()).At(15, 00)
	entries := []log.Entry{
		{chrono.Date(date).At(9, 30), nil, "Designing the timelogger", "SFFEAT0000001"},
		{chrono.Date(date).At(12, 00), nil, "Implementing the timelogger", "SFFEAT0000002"},
		{chrono.Date(date).At(14, 30), nil, "Testing the timelogger", "SFFEAT0000003"},
	}
	actual := Accumulate(entries, date)

	expected := Total([]DayTotal{
		{chrono.DurationOf(5, 30), chrono.GetDay(date), []TaskTotal{
			{chrono.Date(date).At(9, 30), chrono.DurationOf(2, 30), "Designing the timelogger", "SFFEAT0000001", false},
			{chrono.Date(date).At(12, 00), chrono.DurationOf(2, 30), "Implementing the timelogger", "SFFEAT0000002", false},
			{chrono.Date(date).At(14, 30), chrono.DurationOf(0, 30), "Testing the timelogger", "SFFEAT0000003", true},
		}},
	})

	assert.Equal(t, expected, actual)
}
func TestAccumulateWithEOD(t *testing.T) {
	date := chrono.Date(time.Now()).At(15, 00)
	entries := []log.Entry{
		{chrono.Date(date).At(9, 30), nil, "Designing the timelogger", "SFFEAT0000001"},
		{chrono.Date(date).At(12, 00), nil, "Implementing the timelogger", "SFFEAT0000002"},
		{chrono.Date(date).At(14, 30), nil, "Testing the timelogger", "SFFEAT0000003"},
		{chrono.Date(date).At(15, 00), nil, "eod", ""},
	}
	actual := Accumulate(entries, date)

	expected := Total([]DayTotal{
		{chrono.DurationOf(5, 30), chrono.GetDay(date), []TaskTotal{
			{chrono.Date(date).At(9, 30), chrono.DurationOf(2, 30), "Designing the timelogger", "SFFEAT0000001", false},
			{chrono.Date(date).At(12, 00), chrono.DurationOf(2, 30), "Implementing the timelogger", "SFFEAT0000002", false},
			{chrono.Date(date).At(14, 30), chrono.DurationOf(0, 30), "Testing the timelogger", "SFFEAT0000003", false},
			{chrono.Date(date).At(15, 00), chrono.DurationOf(0, 00), "eod", "", false},
		}},
	})

	assert.Equal(t, expected, actual)
}
func TestAccumulateWithEodAndNewDay(t *testing.T) {
	now := chrono.Date(time.Now()).At(15, 00)
	yesterday := now.Add(-24 * time.Hour)
	entries := []log.Entry{
		{chrono.Date(yesterday).At(9, 30), nil, "Designing the timelogger", "SFFEAT0000001"},
		{chrono.Date(yesterday).At(12, 00), nil, "Implementing the timelogger", "SFFEAT0000002"},
		{chrono.Date(yesterday).At(14, 30), nil, "Testing the timelogger", "SFFEAT0000003"},
		{chrono.Date(yesterday).At(15, 00), nil, "eod", ""},
		{chrono.Date(now).At(01, 00), nil, "early start", ""},
	}
	actual := Accumulate(entries, now)

	expected := Total([]DayTotal{
		{chrono.DurationOf(5, 30), chrono.GetDay(yesterday), []TaskTotal{
			{chrono.Date(yesterday).At(9, 30), chrono.DurationOf(2, 30), "Designing the timelogger", "SFFEAT0000001", false},
			{chrono.Date(yesterday).At(12, 00), chrono.DurationOf(2, 30), "Implementing the timelogger", "SFFEAT0000002", false},
			{chrono.Date(yesterday).At(14, 30), chrono.DurationOf(0, 30), "Testing the timelogger", "SFFEAT0000003", false},
			{chrono.Date(yesterday).At(15, 00), chrono.DurationOf(0, 00), "eod", "", false},
		}},
		{chrono.DurationOf(14, 00), chrono.GetDay(now), []TaskTotal{
			{chrono.Date(now).At(1, 00), chrono.DurationOf(14, 00), "early start", "", true},
		}},
	})

	assert.Equal(t, expected, actual)
}
