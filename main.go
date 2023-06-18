package main

import (
	"fmt"
	"github.com/wrigleyster/gorm/util"
	"github.com/wrigleyster/opt"
	"os"
	"strconv"
	"strings"
	"time"
	"wlog/formatter"
	"wlog/log"
	"wlog/manipulation"
)

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
	fmt.Printf("%s: -sid SFFEAT = worked on x [at 9:30] [yesterday|monday-friday]\n", os.Args[0])
	fmt.Printf("%s: -h\n", os.Args[0])
	os.Exit(1)
}
func printLog() {
	db := getDb()
	entries := db.getLogLines(getIntArg(2, 15))
	fmt.Printf("%d entries:\n", len(entries))
	view := formatter.AgendaView(manipulation.Accumulate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func printLogDiff() {
	db := getDb()
	entries := db.getLogLines(getIntArg(2, 15))
	view := formatter.DurationView(manipulation.Aggregate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func printTasks() {
	db := getDb()
	count := getOptionalIntArg(2, 15)
	var tasks []Task
	if count.Exists {
		tasks = db.getTasks(count.Value)
	} else {
		tasks = db.findTasks(strings.Join(os.Args[2:], " "))
	}
	for _, task := range tasks {
		fmt.Printf("%s %s\n", task.ExtId, task.TaskName)
	}
}
func warnOrDie(msg string) {
	print("Warning: " + msg + " Proceed anyway [y/N]: ")
	var response string
	_, err := fmt.Scan(&response)
	util.Log(err)
	response = strings.ToLower(response)
	if response != "y" && response != "yes" {
		os.Exit(1)
	}
}
func getDb() Repository {
	return Seed("sqlite.db")
}
func add() {
	argv := os.Args[1:]
	db := getDb()
	msg := log.Parse(strings.Join(argv, " "))
	if time.Now().Sub(msg.Time) < 0 {
		warnOrDie("That event is in the future.")
	}
	println("add")
	if msg.TaskId != "" {
		task := db.TaskByNameAndExtId(msg.TaskName, msg.TaskId)
		if task.Exists {
			entry := Entry{TaskId: task.Value.Id, StartedAt: msg.Time}
			db.SaveEntry(&entry)
		} else {
			task := Task{ExtId: msg.TaskId, TaskName: msg.TaskName}
			db.SaveTask(&task)
			entry := Entry{TaskId: task.Id, StartedAt: msg.Time}
			db.SaveEntry(&entry)
		}
	} else {
		tasks := db.TaskByName(msg.TaskName)
		if len(tasks) > 1 {
			util.Log("more than one task with that name already exists")
		} else if len(tasks) == 1 {
			task := tasks[0]
			entry := Entry{TaskId: task.Id, StartedAt: msg.Time}
			db.SaveEntry(&entry)
		} else {
			task := Task{TaskName: msg.TaskName}
			db.SaveTask(&task)
			entry := Entry{TaskId: task.Id, StartedAt: msg.Time}
			db.SaveEntry(&entry)
		}
	}
}
func setId() {
	println("Not Implemented yet")
}
func deleteEntry() {
	println("Not Implemented yet")
}
func parseArgs(argv []string) func() {
	if len(argv) == 0 ||
		argv[0] == "-l" {
		return printLog
	} else if argv[0] == "-ld" {
		return printLogDiff
	} else if argv[0] == "-lt" {
		return printTasks
	} else if argv[0] == "-sid" {
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
