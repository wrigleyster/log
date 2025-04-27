package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"wlog/chrono"
	"wlog/formatter"
	"wlog/log"
	"wlog/manipulation"
	"wlog/model"

	"github.com/wrigleyster/gorm/util"
	"github.com/wrigleyster/opt"
)

func getDb() model.Repository {
	db := model.Seed("sqlite.db")
	db.CleanChildlessParents()
	return db
}
func getArg(i int, fallback string) string {
	if len(os.Args) > i {
		return os.Args[i]
	}
	return fallback
}
func getIntArg(i, fallback int) int {
	if len(os.Args) > i {
		if i, e := strconv.Atoi(os.Args[i]); e == nil {
			return i
		}
	}
	return fallback
}
func getOptionalIntArg(i, fallback int) opt.Maybe[int] {
	if len(os.Args) > i {
		if i, e := strconv.Atoi(os.Args[i]); e == nil {
			return opt.Some(i)
		}
		return opt.No[int]()
	}
	return opt.Some(fallback)
}

func printUsage() {
	fmt.Printf("%s: [SFFEAT] working on x [at 9:30] [yesterday|monday-friday]\n", os.Args[0])
	fmt.Printf("%s: -l[d|t] [count]\n", os.Args[0])
	fmt.Printf("%s: -dd [SFFEAT] worked on x at 9:30 [yesterday|monday-friday]\n", os.Args[0])
	fmt.Printf("%s: -s SFFEAT = worked on x [at 9:30] [yesterday|monday-friday]\n", os.Args[0])
	fmt.Printf("%s: -h\n", os.Args[0])
	os.Exit(1)
}
func printLog() {
	db := getDb()
	entries := db.GetLogLines(getIntArg(2, 15))
	fmt.Printf("%d entries:\n", len(entries))
	view := formatter.AgendaView(manipulation.Accumulate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func printLogDiff() {
	db := getDb()
	entries := db.GetLogLines(getIntArg(2, 15))
	view := formatter.DurationView(manipulation.Aggregate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func printDailyLog() {
	db := getDb()
	words := strings.Split(getArg(2, ""), " ")
	date := chrono.ParseDate(words[0], time.Now())
	entries := db.GetDailyLog(date)
	view := formatter.AgendaView(manipulation.Accumulate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func printDailyLogDiff() {
	db := getDb()
	words := strings.Split(getArg(2, ""), " ")
	date := chrono.ParseDate(words[0], time.Now())
	entries := db.GetDailyLog(date)
	view := formatter.DurationView(manipulation.Aggregate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func printTasks() {
	db := getDb()
	count := getOptionalIntArg(2, 15)
	var tasks []model.Task
	if count.Exists {
		tasks = db.GetTasks(count.Value)
	} else {
		tasks = db.FindTasks(strings.Join(os.Args[2:], " "))
	}
	for _, task := range tasks {
		fmt.Printf("%s %s\n", task.TaskName, task.ExtId)
	}
}
func Prompt(prompt ...string) string {
	for _, msg := range prompt {
		print(msg)
	}
	var reply string
	_, err := fmt.Scanln(&reply)
	util.Log(err)
	return reply
}
func warnOrDie(msg string) {
	response := Prompt("Warning:", msg, "Proceed anyway [y/N]: ")
	response = strings.ToLower(response)
	if response != "y" && response != "yes" {
		os.Exit(1)
	}
}
func add() {
	argv := os.Args[1:]
	db := getDb()
	msg := log.Parse(argv)
	if time.Now().Sub(msg.Time) < 0 {
		warnOrDie("That event is in the future.")
	}
	println("add")
	if msg.TaskId != "" {
		if task := db.TaskByNameAndExtId(msg.TaskName, msg.TaskId); task.Exists {
			entry := model.Entry{TaskId: task.Value.Id, StartedAt: msg.Time}
			db.SaveEntry(&entry)
		} else {
			task := model.Task{ExtId: msg.TaskId, TaskName: msg.TaskName}
			db.SaveTask(&task)
			entry := model.Entry{TaskId: task.Id, StartedAt: msg.Time}
			db.SaveEntry(&entry)
		}
	} else {
		tasks := db.TasksByName(msg.TaskName)
		if len(tasks) > 1 {
			util.Log("more than one task with that name already exists")
			return
		}
		var task model.Task
		if len(tasks) == 1 {
			task = tasks[0]
		} else {
			task = model.Task{TaskName: msg.TaskName}
		}
		db.SaveTask(&task)
		entry := model.Entry{TaskId: task.Id, StartedAt: msg.Time}
		db.SaveEntry(&entry)
	}
}
func setId() {
	argv := os.Args[1:]
	db := getDb()
	extId, name := log.ParseSet(argv)
	if extId != "" && name != "" {
		tasks := db.TasksByName(name)
		if len(tasks) == 1 {
			tasks[0].ExtId = extId
			db.SaveTask(&tasks[0])
		} else {
			println("Error: multiple tasks named: ", name)
			os.Exit(1)
		}
	} else {
		println("Error: invalid assignment")
		os.Exit(1)
	}
}
func deleteEntry() {
	argv := os.Args[2:]
	db := getDb()
	msg := log.Parse(argv)
	var task opt.Maybe[model.Task]
	if entry := db.EntryByTimestamp(msg.Time); entry.Exists {
		if task := db.TaskById(entry.Value.TaskId); task.Exists {
			prompt := fmt.Sprintf("Would you like to delete \"%s %s %s\" [y/N]: ", entry.Value.StartedAt, task.Value.TaskName, task.Value.ExtId)
			if reply := strings.ToLower(Prompt(prompt)); reply == "y" {
				db.DeleteEntry(entry.Value)
			}
			return
		}
	}
	println("trying generic")
	if msg.TaskId != "" {
		task = db.TaskByNameAndExtId(msg.TaskName, msg.TaskId)
	} else {
		task = opt.First(db.TasksByName(msg.TaskName)) // TODO pick in sort order
	}
	if task.Exists {
		if entry := opt.First(db.EntriesByTaskId(task.Value.Id)); entry.Exists {
			prompt := fmt.Sprintf("Would you like to delete \"%s %s %s\" [y/N]: ", entry.Value.StartedAt, task.Value.TaskName, task.Value.ExtId)
			if reply := strings.ToLower(Prompt(prompt)); reply == "y" {
				db.DeleteEntry(entry.Value)
			}
			return
		}
	}
	fmt.Println("unable to find entry.")
}
func parseArgs(argv []string) func() {
	if len(argv) == 0 ||
		argv[0] == "-l" {
		return printLog
	} else if argv[0] == "-ld" {
		return printLogDiff
	} else if argv[0] == "-ll" {
		return printDailyLog
	} else if argv[0] == "-lld" {
		return printDailyLogDiff
	} else if argv[0] == "-lt" {
		return printTasks
	} else if argv[0] == "-s" {
		return setId
	} else if argv[0] == "-dd" {
		return deleteEntry
	} else if argv[0] == "-h" {
		return printUsage
	} else if argv[0] == os.Args[0] {
		return printUsage
	} else {
		return add
	}
}

func main() {
	argv := os.Args[1:]
	action := parseArgs(argv)
	action()
}
func M[T, U any](n opt.Maybe[T], f func(T) U) opt.Maybe[U] {
	if n.Exists {
		return opt.Some(f(n.Value))
	}
	return opt.No[U]()
}

type Row interface {
	Next() bool
	Scan(dest ...any) error
	Columns() ([]string, error)
}
