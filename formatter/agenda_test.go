package formatter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"wlog/chrono"
	"wlog/manipulation"
	"wlog/model"
)

func TestFormat(t *testing.T) {
	date := time.Date(2023, 4, 23, 14, 37, 0, 0, time.Local)
	entries := []model.LogEntry{
		{Time: chrono.Date(date).At(9, 30), TaskName: "Designing the timelogger", ExtId: "SFFEAT0000001"},
		{Time: chrono.Date(date).At(12, 00), TaskName: "Implementing the timelogger", ExtId: "SFFEAT0000002"},
		{Time: chrono.Date(date).At(14, 30), TaskName: "Testing the timelogger", ExtId: "SFFEAT0000003"},
	}
	actual := AgendaView(manipulation.Accumulate(entries, date)).Format(Ascending)

	expected := "Sunday 23 April\n" +
		" 09:30 SFFEAT0000001 Designing the timelogger\n" +
		" 12:00 SFFEAT0000002 Implementing the timelogger\n" +
		" 14:30 SFFEAT0000003 Testing the timelogger"

	assert.Equal(t, expected, actual)
}
