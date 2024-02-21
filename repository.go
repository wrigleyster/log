package main

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/wrigleyster/gorm"
	"github.com/wrigleyster/gorm/sqlite"
	"github.com/wrigleyster/gorm/util"
	"github.com/wrigleyster/opt"
	"time"
	"wlog/list"
	"wlog/log"
)

type Repository struct {
	db gorm.DataSource
}

func Seed(name string) Repository {
	db := sqlite.New(name)
	db.With("id text not null, extId text, taskName text not null, state text").
		Key("id").
		CreateTable("task")
	db.With("id text not null, taskId text not null, startedAt timestamp unique not null").
		Key("id").
		CreateTable("entry")
	return Repository{db}
}

type Task struct {
	Id, ExtId, TaskName, State string
}
type Entry struct {
	Id, TaskId      string
	StartedAt       time.Time
}

func (task Task) fields() []any {
	return []any{task.Id, task.ExtId, task.TaskName, task.State}
}
func (entry Entry) fields() []any {
	return []any{entry.Id, entry.TaskId, entry.StartedAt}
}
func makeTask(row *sql.Rows) Task {
	var task Task
	err := row.Scan(&task.Id, &task.ExtId, &task.TaskName, &task.State)
	util.Log(err)
	return task
}
func makeEntry(row *sql.Rows) Entry {
	var entry Entry
	err := row.Scan(&entry.Id, &entry.TaskId, &entry.StartedAt)
	util.Log(err)
	return entry
}
func (repo Repository) SaveTask(task *Task) {
	if task.Id == "" {
		task.Id = uuid.NewString()
	}
	repo.db.From("task").
		Replace(task.fields()...)
}
func (repo Repository) SaveEntry(entry *Entry) {
	if entry.Id == "" {
		entry.Id = uuid.NewString()
	}
	repo.db.From("entry").
		Replace(entry.fields()...)
}
func (repo Repository) Save(entry log.Entry) {

	res := repo.db.From("task").
		Where("extId = ? or taskName = ?", entry.TaskId, entry.TaskName).
		Select("id")
	res.Next()
	taskid := uuid.NewString()
	uuid.New().ID()
	repo.db.From("task").Replace(taskid, entry.TaskId, entry.TaskName,"")
	repo.db.From("entry").Replace(uuid.NewString(), taskid, entry.Time.String())

}
func (repo Repository) TaskById(id string) opt.Maybe[Task] {
	return opt.First(repo.taskBy("id = ?", id))
}
func (repo Repository) TasksByExtId(extId string) []Task {
	return repo.taskBy("extId = ?", extId)
}
func (repo Repository) TasksByName(name string) []Task {
	return repo.taskBy("taskName = ?", name)
}
func (repo Repository) TaskByNameAndExtId(name, extId string) opt.Maybe[Task] {
	return opt.First(repo.taskBy("taskName = ? and extId = ?", name, extId))
}
func (repo Repository) EntryById(id string) opt.Maybe[Entry] {
	return opt.First(repo.entryBy("id = ?", id))
}
func (repo Repository) EntryByTimestamp(startedAt time.Time) opt.Maybe[Entry] {
	return opt.First(repo.entryBy("startedAt = ?", startedAt))
}
func (repo Repository) EntriesByTaskId(taskId string) []Entry {
	return repo.entryBy("taskId = ?", taskId)
}

func (repo Repository) taskBy(predicate string, value ...any) []Task {
	res := repo.db.From("task").
		Where(predicate, value...).
		Select()
	var tasks []Task
	for res.Next() {
		tasks = append(tasks, makeTask(res))
	}
	return tasks
}
func (repo Repository) entryBy(predicate string, value ...any) []Entry {
	res := repo.db.From("entry").
		Where(predicate, value...).
		Select()
	var entries []Entry
	for res.Next() {
		entries = append(entries, makeEntry(res))
	}
	return entries
}
func (repo Repository) getLogLines(count int) []log.Entry {
	var entries []log.Entry
	repo.db.Orm(func(db *sql.DB) {
		stmt, err := db.Prepare("SELECT startedAt, task.taskName, task.extId FROM task INNER JOIN entry ON task.id = entry.taskId ORDER BY entry.startedAt DESC limit ?")
		util.Log(err)
		row, err := stmt.Query(count)
		util.Log(err)
		for row.Next() {
			entry := log.Entry{}
			err = row.Scan(&entry.Time, &entry.TaskName, &entry.TaskId)
			util.Log(err)
			entries = append(entries, entry)
		}
	})
	list.Reverse(entries)
	return entries
}
func (repo Repository) getTasks(count int) []Task {
	var tasks []Task
	repo.db.Orm(func(db *sql.DB) {
		stmt, err := db.Prepare("SELECT * FROM task WHERE state != 'done' LIMIT ?")
		util.Log(err)
		row, err := stmt.Query(count)
		util.Log(err)
		for row.Next() {
			tasks = append(tasks, makeTask(row))
		}
	})
	return tasks
}
func (repo Repository) findTasks(nameOrExtId string) []Task {
    query := "%" + nameOrExtId + "%"
    return repo.taskBy("state != 'done' and taskName like ? or extId like ?", query, query)
}
func (repo Repository) CleanChildlessParents() (rowsCleaned int64) {
	repo.db.Orm(func(db *sql.DB) {
		stmt, err := db.Exec("DELETE FROM task WHERE task.id IN (SELECT task.id FROM task LEFT OUTER JOIN entry ON task.id = entry.taskId WHERE entry.taskId IS NULL)")
		util.Log(err)
		rowsCleaned, err = stmt.RowsAffected()
		util.Log(err)
	})
	return rowsCleaned
}
func (repo Repository) DeleteEntry(entry Entry){
	repo.db.From("entry").Where("id = ?", entry.Id).Delete()
}

//func (repo Repository) ByKey(key string) (Entry, error) {
//	res := repo.db.From(repo.table).
//		Where("key = ?", key).
//		Select()
//	entry := Entry{}
//	for res.Next() {
//		return entry.fill(res), nil
//	}
//	return entry, errors.New("repository: key not found")
//}
