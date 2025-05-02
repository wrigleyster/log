package action

import (
	"fmt"
	"github.com/wrigleyster/gorm/util"
	"time"
	"wlog/log"
	"wlog/model"
)

func UsageAdd(argv Argv) {
	fmt.Printf("%s: [SFFEAT] working on x [at 9:30] [yesterday|monday-friday]\n", argv[0])
}

func Add(db *model.Repository, argv Argv) {
	msg := log.Parse(argv)
	if time.Now().Sub(msg.Time) < 0 {
		warnOrDie("That event is in the future.")
	}
	println("add")
	if msg.ExtId != "" {
		if task := db.TaskByNameAndExtId(msg.TaskName, msg.ExtId); task.Exists {
			entry := model.Entry{TaskId: task.Value.Id, StartedAt: msg.Time}
			db.SaveEntry(&entry)
		} else {
			task := model.Task{ExtId: msg.ExtId, TaskName: msg.TaskName}
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
