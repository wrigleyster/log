package main

import (
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
	}
	repo.SaveEntry(&entry)
	t.Logf("entry = %s", entry.Id)
	actual := repo.EntryById(entry.Id)
	assert.True(t, actual.Exists)
	assert.Equal(t, entry.TaskId, actual.Value.TaskId)
}

func TestCleanupChildlessParents(t *testing.T) {
	dbname := "test.db"
	t.Cleanup(func() {
		err := os.Remove(dbname)
		util.Log(err)
	})
	repo := Seed(dbname)

	task := Task{
		Id: "",
		TaskName: "A",
	}
	repo.SaveTask(&task)
	entry := Entry{
		Id: "",
		TaskId: task.Id,
	}
	repo.SaveEntry(&entry)
	childless := Task{
		TaskName: "B",
	}
	repo.SaveTask(&childless)


	assert.Equal(t, 2, len(repo.getTasks(10)))
	repo.CleanChildlessParents()
	assert.Equal(t, 1, len(repo.getTasks(10)))
	repo.DeleteEntry(entry)
	repo.CleanChildlessParents()
	assert.Equal(t, 0, len(repo.getTasks(10)))

}
