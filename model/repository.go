package model

import (
	"database/sql"
	"github.com/wrigleyster/gorm/sqlite"
	"time"
	"wlog/chrono"
	"wlog/list"

	"github.com/google/uuid"
	"github.com/wrigleyster/gorm"
	"github.com/wrigleyster/gorm/util"
	"github.com/wrigleyster/opt"
)

type Repository struct {
	db      gorm.DataSource
	IsLocal bool
}

func Local(db gorm.DataSource) Repository {
	return Repository{db, true}
}
func SqLite(name string) Repository {
	return Local(sqlite.New(name))
}
func Remote(db gorm.DataSource) Repository {
	return Repository{db, false}
}

func (repo *Repository) Seed() {
	repo.db.With("id text not null, extId text, taskName text not null, state text").
		Key("id").
		CreateTable("task")
	repo.db.With("id text not null, taskId text not null, startedAt timestamp unique not null").
		Key("id").
		CreateTable("entry")
}

type Table interface {
	id() string
	setId(string)
	fields() []any
	tableName() string
}
type Factory[T any] interface {
	Make(row *sql.Rows) T
	Table() string
}

type Task struct {
	Id, ExtId, TaskName, State string
}
type Entry struct {
	Id, TaskId string
	StartedAt  time.Time
}
type LogEntry struct {
	Time     time.Time
	TaskName string
	ExtId    string
}

func (task Task) fields() []any {
	return []any{task.Id, task.ExtId, task.TaskName, task.State}
}
func (entry Entry) fields() []any {
	return []any{entry.Id, entry.TaskId, entry.StartedAt}
}
func (Task) tableName() string {
	return TaskFactory{}.Table()
}
func (Entry) tableName() string {
	return EntryFactory{}.Table()
}

type EntryFactory struct{}

func (EntryFactory) Table() string {
	return "entry"
}
func (EntryFactory) Make(row *sql.Rows) Entry {
	var entry Entry
	err := row.Scan(&entry.Id, &entry.TaskId, &entry.StartedAt)
	util.Log(err)
	return entry
}

type TaskFactory struct{}

func (TaskFactory) Table() string {
	return "task"
}
func (TaskFactory) Make(row *sql.Rows) Task {
	var task Task
	err := row.Scan(&task.Id, &task.ExtId, &task.TaskName, &task.State)
	util.Log(err)
	return task
}
func makeTask(row *sql.Rows) Task {
	return TaskFactory{}.Make(row)
}
func makeEntry(row *sql.Rows) Entry {
	return EntryFactory{}.Make(row)
}
func (repo *Repository) save(table Table) {
	if table.id() == "" {
		table.setId(uuid.NewString())
	}
	repo.db.From(table.tableName()).
		Replace(table.fields()...)
}
func (repo *Repository) SaveTask(task *Task) {
	if task.Id == "" {
		task.Id = uuid.NewString()
	}
	repo.db.From("task").
		Replace(task.fields()...)
}
func (repo *Repository) SaveEntry(entry *Entry) {
	if entry.Id == "" {
		entry.Id = uuid.NewString()
	}
	repo.db.From("entry").
		Replace(entry.fields()...)
}
func (repo *Repository) Save(entry LogEntry) {

	res := repo.db.From("task").
		Where("extId = ? or taskName = ?", entry.ExtId, entry.TaskName).
		Select("id")
	var taskid string
	if res.Next() {
		err := res.Scan(&taskid)
		util.Log(err)
		if res.Next() {
			panic("duplicate task")
		}
		err = res.Close()
		util.Log(err)
	} else {
		taskid = uuid.NewString()
	}

	repo.db.From("task").Replace(taskid, entry.ExtId, entry.TaskName, "")
	repo.db.From("entry").Replace(uuid.NewString(), taskid, entry.Time.String())

}
func (repo *Repository) TaskById(id string) opt.Maybe[Task] {
	return opt.First(repo.taskBy("id = ?", id))
}
func (repo *Repository) TasksByExtId(extId string) []Task {
	return repo.taskBy("extId = ?", extId)
}
func (repo *Repository) TasksByName(name string) []Task {
	return repo.taskBy("taskName = ?", name)
}
func (repo *Repository) TaskByNameAndExtId(name, extId string) opt.Maybe[Task] {
	return opt.First(repo.taskBy("taskName = ? and extId = ?", name, extId))
}
func (repo *Repository) EntryById(id string) opt.Maybe[Entry] {
	return opt.First(repo.entryBy("id = ?", id))
}
func (repo *Repository) EntryByTimestamp(startedAt time.Time) opt.Maybe[Entry] {
	return opt.First(repo.entryBy("startedAt = ?", startedAt))
}
func (repo *Repository) EntriesByTaskId(taskId string) []Entry {
	return repo.entryBy("taskId = ?", taskId)
}
func findBy[T any](repo *Repository, factory Factory[T], predicate string, predicateValues ...any) []T {
	res := repo.db.From(factory.Table()).
		Where(predicate, predicateValues...).
		Select()
	defer res.Close()
	var entries []T
	for res.Next() {
		entries = append(entries, factory.Make(res))
	}
	return entries
}
func (repo *Repository) taskBy(predicate string, values ...any) []Task {
	return findBy(repo, TaskFactory{}, predicate, values...)
}
func (repo *Repository) entryBy(predicate string, values ...any) []Entry {
	return findBy(repo, EntryFactory{}, predicate, values...)
}
func (repo *Repository) GetLogLines(count int) []LogEntry {
	var entries []LogEntry
	repo.db.Orm(func(db *sql.DB) {
		stmt, err := db.Prepare("SELECT startedAt, task.taskName, task.extId FROM task INNER JOIN entry ON task.id = entry.taskId ORDER BY entry.startedAt DESC limit ?")
		util.Log(err)
		row, err := stmt.Query(count)
		util.Log(err)
		for row.Next() {
			entry := LogEntry{}
			err = row.Scan(&entry.Time, &entry.TaskName, &entry.ExtId)
			util.Log(err)
			entries = append(entries, entry)
		}
	})
	list.Reverse(entries)
	return entries
}
func (repo *Repository) GetDailyLog(date time.Time) []LogEntry {
	var entries []LogEntry
	date = chrono.Date(date).At(0, 0)
	row := repo.db.From("task").
		InnerJoin("entry e", "task.id = e.taskId").
		Where("? < e.startedAt and e.startedAt < ?", date, date.AddDate(0,0,1)).
		Select("e.startedAt", "task.taskName", "task.extId")
	defer row.Close()
	for row.Next() {
			entry := LogEntry{}
			err := row.Scan(&entry.Time, &entry.TaskName, &entry.ExtId)
			util.Log(err)
			entries = append(entries, entry)
	}
	return entries
}
func (repo *Repository) GetTasks(count int) []Task {
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
func (repo *Repository) FindTasks(nameOrExtId string) []Task {
	query := "%" + nameOrExtId + "%"
	return repo.taskBy("state != 'done' and taskName like ? or extId like ?", query, query)
}
func (repo *Repository) CleanChildlessParents() (rowsCleaned int64) {
	repo.db.Orm(func(db *sql.DB) {
		stmt, err := db.Exec("DELETE FROM task WHERE task.id IN (SELECT task.id FROM task LEFT OUTER JOIN entry ON task.id = entry.taskId WHERE entry.taskId IS NULL)")
		util.Log(err)
		rowsCleaned, err = stmt.RowsAffected()
		util.Log(err)
	})
	return rowsCleaned
}
func (repo *Repository) DeleteEntry(entry Entry) {
	repo.db.From("entry").Where("id = ?", entry.Id).Delete()
}

//func (repo *Repository) ByKey(key string) (Entry, error) {
//	res := repo.db.From(repo.table).
//		Where("key = ?", key).
//		Select()
//	entry := Entry{}
//	for res.Next() {
//		return entry.fill(res), nil
//	}
//	return entry, errors.New("repository: key not found")
//}
