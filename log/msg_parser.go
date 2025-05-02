package log

import (
	"fmt"
	"github.com/wrigleyster/gorm/util"
	"strings"
	"time"
	"wlog/chrono"
	"wlog/list"
	"wlog/model"
)

type Entry model.LogEntry

func NewLogEntry(input string) Entry {
	return Entry{time.Now(), input, ""}
}
func (entry Entry) ParseTime() Entry {
	now := entry.Time.Truncate(time.Minute)
	words := strings.Split(entry.TaskName, " ")
	if len(words) > 1 && strings.Contains(words[0], ":") && !strings.HasSuffix(words[0], ":"){
		startTime := chrono.ParseOptionalTime(words[0], entry.Time)
		if !startTime.Exists {
			return entry
		}
		entry.Time = startTime.Value
		entry.TaskName = strings.Join(words[1:], " ")
	} else if len(words) > 2 && words[len(words)-2] == "at" {
		startTime := chrono.ParseOptionalTime(list.El(words, -1), entry.Time)
		if !startTime.Exists {
			return entry
		}
		entry.Time = startTime.Value
		entry.TaskName = strings.Join(list.Sl(words, 0, -2), " ")
	} else {
		entry.Time = now
	}
	return entry

}
func (entry Entry) parseExtId() Entry {
	words := strings.Split(entry.TaskName, " ")
	for i, v := range words {
		if strings.HasPrefix(v, "SFFEAT") ||
			strings.HasPrefix(v, "SFSTRY") {
			entry.ExtId = v
			entry.TaskName = strings.Join(append(words[:i], words[i+1:]...), " ")
			break
		}

	}
	return entry
}

func (entry Entry) parseDate() Entry {
	words := strings.Split(entry.TaskName, " ")
	newDate := chrono.RelativeDate(entry.Time, list.El(words, -1))
	if newDate == entry.Time {
		return entry
	}
	entry.TaskName = strings.Join(list.Sl(words, 0, -1), " ")
	entry.Time = newDate
	return entry
}

func (entry Entry) parseFrontDate() Entry {
	words := strings.Split(entry.TaskName, " ")
	newDate := chrono.ParseDate(words[0], entry.Time)
	if newDate == entry.Time {
		return entry
	}
	entry.TaskName = strings.Join(words[1:], " ")
	entry.Time = newDate
	return entry
}

func (entry Entry) IsEOD() bool {
	return strings.ToLower(entry.TaskName) == "eod"
}

func (entry Entry) Str() string {
	return fmt.Sprintf("Entry(%s,%s,%s)", entry.Time.String(), entry.TaskName, entry.ExtId)
}

func Parse(argv []string) Entry {
	input := strings.Join(argv, " ")
	entry := NewLogEntry(input).
		parseDate().
		parseFrontDate().
		ParseTime().
		parseExtId()
	return entry
}

func Add(db *model.Repository, msg Entry) {
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
