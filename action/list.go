package action

import (
	"fmt"
	"os"
	"strings"
	"time"
	"wlog/chrono"
	"wlog/formatter"
	"wlog/manipulation"
	"wlog/model"
)

func UsageList(argv Argv) {
	fmt.Printf("%s: -l[d|t] [count]\n", argv[0])
	fmt.Printf("%s: -ll[d] [count]\n", argv[0])
}

func ListLog(db *model.Repository, argv Argv) {
	entries := db.GetLogLines(argv.getIntArg(2, 15))
	fmt.Printf("%d entries:\n", len(entries))
	view := formatter.AgendaView(manipulation.Accumulate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func ListLogDiff(db *model.Repository, argv Argv) {
	entries := db.GetLogLines(argv.getIntArg(2, 15))
	view := formatter.DurationView(manipulation.Aggregate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func ListDailyLog(db *model.Repository, argv Argv) {
	words := strings.Split(argv.getArg(2, ""), " ")
	date := chrono.ParseDate(words[0], time.Now())
	entries := db.GetDailyLog(date)
	if entries == nil && !chrono.IsWeekday(date) {
		die("it's the weekend on " + chrono.Date(date).Iso())
	} else if entries == nil {
		die("no data for " + chrono.Date(date).Iso())
	}
	view := formatter.AgendaView(manipulation.Accumulate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func ListDailyLogDiff(db *model.Repository, argv Argv) {
	words := strings.Split(argv.getArg(2, ""), " ")
	date := chrono.ParseDate(words[0], time.Now())
	entries := db.GetDailyLog(date)
	if entries == nil && !chrono.IsWeekday(date) {
		die("it's the weekend on " + chrono.Date(date).Iso())
	} else if entries == nil {
		die("no data for " + chrono.Date(date).Iso())
	}
	view := formatter.DurationView(manipulation.Aggregate(entries, time.Now()))
	println(view.Format(formatter.Ascending))
}
func ListTasks(db *model.Repository, argv Argv) {
	count := argv.getOptionalIntArg(2, 15)
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
