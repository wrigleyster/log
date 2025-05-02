package verifier

import (
	"github.com/wrigleyster/gorm/util"
	"os"
	"testing"
	"time"
	"wlog/model"

	"github.com/stretchr/testify/assert"
)

func TestWeekday(t *testing.T) {
	dbname := "test.db"
	t.Cleanup(func() {
		err := os.Remove(dbname)
		util.Log(err)
	})
	repo := model.Seed(dbname)

	v := New(&repo)
	monday := time.Date(2025, time.January, 6, 1, 1, 1, 1, time.UTC)
	tuesday := monday.AddDate(0, 0, 1)
	wednesday := monday.AddDate(0, 0, 2)
	thursday := monday.AddDate(0, 0, 3)
	friday := monday.AddDate(0, 0, 4)
	saturday := monday.AddDate(0, 0, 5)
	sunday := monday.AddDate(0, 0, 6)

	assert.True(t, v.weekday(monday))
	assert.True(t, v.weekday(tuesday))
	assert.True(t, v.weekday(wednesday))
	assert.True(t, v.weekday(thursday))
	assert.True(t, v.weekday(friday))

	assert.False(t, v.weekday(saturday))
	assert.False(t, v.weekday(sunday))
}

func TestBegun(t *testing.T) {
	dbname := "test.db"
	t.Cleanup(func() {
		err := os.Remove(dbname)
		util.Log(err)
	})
	repo := model.Seed(dbname)
	v := New(&repo)
	monday := time.Date(2025, time.January, 6, 1, 1, 1, 1, time.UTC)

	assert.False(t, v.began(monday))
	v.Repo.Save(model.LogEntry{Time: monday, TaskName: "Designing the timelogger", ExtId: "SFFEAT0000001"})
	assert.True(t, v.began(monday))

	assert.False(t, v.lunch(monday))
	v.Repo.Save(model.LogEntry{Time: monday.Add(12 * time.Hour), TaskName: "lunch", ExtId: ""})
	assert.True(t, v.lunch(monday))

	assert.False(t, v.eod(monday))
	v.Repo.Save(model.LogEntry{Time: monday.Add(14 * time.Hour), TaskName: "eod", ExtId: ""})
	assert.True(t, v.eod(monday))
}
