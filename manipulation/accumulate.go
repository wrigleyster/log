package manipulation

import (
	"fmt"
	"os"
	"strings"
	"time"
	"wlog/chrono"
	"wlog/log"
)

type Total []DayTotal

type DayTotal struct {
	Duration chrono.Duration
	Day      chrono.Day
	Tasks    []TaskTotal // todo replace with aggregating map
}

func NewDayTotal(t time.Time) DayTotal {
	return DayTotal{Day: chrono.GetDay(t)}
}

type TaskTotal struct {
	Duration    chrono.Duration
	Name, ExtId string
	IsOpen      bool
}

func ternary[C any](predicate bool, t, f C) C {
	if predicate {
		return t
	}
	return f
}
func (t TaskTotal) Str() string {
	openness := ternary(t.IsOpen, "+", " ")
	if t.ExtId == "" {
		return fmt.Sprintf("%s%s %s", openness, t.Duration.Str(), t.Name)
	}
	return fmt.Sprintf("%s%s %s %s", openness, t.Duration.Str(), t.ExtId, t.Name)
}
func (t TaskTotal) IsEOD() bool {
	return strings.ToLower(t.Name) == "eod"
}

func getTaskTotal(task log.Entry, endTime time.Time, isOpen bool) TaskTotal {
	if task.IsEOD() {
		endTime = task.Time
		isOpen = false
	}
	return TaskTotal{
		chrono.GetDuration(task.Time, endTime),
		task.TaskName,
		task.TaskId,
		isOpen,
	}
}
func last[T any](slice []T) int {
	return len(slice) - 1
}
func assertAscending(entries []log.Entry) {
	if len(entries) < 2 {
		return
	}
	if entries[0].Time.Sub(entries[1].Time) > 0 {
		println("Assert: Expected entries to be in ascending order")
		os.Exit(1)
	}
	assertAscending(entries[1:])
}
func Accumulate(entries []log.Entry, now time.Time) Total {
	assertAscending(entries)
	total := Total{}
	dayTotal := DayTotal{}
	var task TaskTotal
	for i, j := 0, 1; j < len(entries); i, j = j, j+1 {
		entry := entries[i]
		endTime := entries[j].Time
		if i == 0 {
			dayTotal = NewDayTotal(entry.Time)
		}
		task = getTaskTotal(entry, endTime, false)

		dayTotal.Duration = dayTotal.Duration.Add(task.Duration)
		dayTotal.Tasks = append(dayTotal.Tasks, task)
		if chrono.GetDay(entry.Time) != chrono.GetDay(endTime) {
			total = append(total, dayTotal)
			dayTotal = NewDayTotal(endTime)
		}
	}
	entry := entries[last(entries)]
	task = getTaskTotal(entry, now, true)
	dayTotal.Duration = dayTotal.Duration.Add(task.Duration)
	dayTotal.Tasks = append(dayTotal.Tasks, task)

	total = append(total, dayTotal)
	return total
}
