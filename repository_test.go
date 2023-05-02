package main

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/wrigleyster/gorm/util"
	"os"
	"testing"
	"time"
)

func TestSaveAndRecover(t *testing.T) {
	dbname := "test.db"
	t.Cleanup(func() {
		err := os.Remove(dbname)
		util.Log(err)
	})
	repo := Seed(dbname)

	entry := Entry{
		Id:              "",
		TaskId:          "SFFEAT012345",
		StartedAt:       time.Now(),
		DurationMinutes: sql.NullInt32{0, false},
	}
	repo.SaveEntry(&entry)
	t.Logf("entry = %s", entry.Id)
	actual := repo.EntryById(entry.Id)
	assert.True(t, actual.Exists)
	assert.Equal(t, entry.TaskId, actual.Value.TaskId)
}
