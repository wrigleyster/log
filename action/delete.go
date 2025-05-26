package action

import (
	"fmt"
	"github.com/wrigleyster/opt"
	"strings"
	"wlog/log"
	"wlog/model"
)

func UsageDelete(argv Argv) {
	fmt.Println("Delete")
	fmt.Printf("\t%s: -dd [SFFEAT] worked on x at 9:30 [yesterday|monday-friday]\n", argv[0])
}

func DeleteEntry(db *model.Repository, argv Argv) {
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
	if msg.ExtId != "" {
		task = db.TaskByNameAndExtId(msg.TaskName, msg.ExtId)
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
