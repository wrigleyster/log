package main

import (
	"fmt"
	"github.com/wrigleyster/gorm/util"
	"github.com/wrigleyster/opt"
	"os"
	"strings"
	"time"
	"wlog/formatter"
	"wlog/log"
)

func printUsage() {
	fmt.Printf("%s: [SFFEAT] working on x [at 9:30] [yesterday|monday-friday]\n", os.Args[0])
	fmt.Printf("%s: -h\n", os.Args[0])
	os.Exit(1)
}
func printLog() {
	db := Seed("sqlite.db")
	entries := db.getLogLines(15)
	fmt.Printf("%d entries:\n", len(entries))
	println(formatter.Format(entries))
}
func printLogDiff() {
	db := Seed("sqlite.db")
	entries := db.getLogLines(15)
	println(formatter.FormatDurations(entries, time.Now()))
}
func add() {
	println("add")
	argv := os.Args[1:]
	db := Seed("sqlite.db")
	msg := log.Parse(strings.Join(argv, " "))
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
func parseArgs(argv []string) func() {
	if len(argv) == 0 ||
		argv[0] == "-l" {
		return printLog
	} else if argv[0] == "-ld" {
		return printLogDiff
	} else if argv[0] == "-h" {
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
