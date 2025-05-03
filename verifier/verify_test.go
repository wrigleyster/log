package verifier

import (
	"github.com/wrigleyster/gorm/util"
	"os"
	"testing"
	"time"
	"wlog/model"

	"github.com/stretchr/testify/assert"
)

func TestBegun(t *testing.T) {
	dbname := "test.db"
	t.Cleanup(func() {
		err := os.Remove(dbname)
		util.Log(err)
	})
	repo := model.SqLite(dbname)
	repo.Seed()
	monday := time.Date(2025, time.January, 6, 1, 1, 1, 1, time.UTC)
	v := New(&repo, monday)

	assert.False(t, v.began)
	repo.Save(model.LogEntry{Time: monday, TaskName: "Designing the timelogger", ExtId: "SFFEAT0000001"})
	v = New(&repo, monday)
	assert.True(t, v.began)

	assert.False(t, v.lunch)
	repo.Save(model.LogEntry{Time: monday.Add(12 * time.Hour), TaskName: "lunch", ExtId: ""})
	v = New(&repo, monday)
	assert.True(t, v.lunch)

	assert.False(t, v.eod)
	repo.Save(model.LogEntry{Time: monday.Add(14 * time.Hour), TaskName: "eod", ExtId: ""})
	v = New(&repo, monday)
	assert.True(t, v.eod)
}
