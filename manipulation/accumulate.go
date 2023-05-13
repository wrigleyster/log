package manipulation

import (
	"fmt"
	"os"
	"time"
	"wlog/log"
)

type Total []DayTotal
type Day time.Time

func GetDay(t time.Time) Day {
	return NewDay(t.Year(), t.Month(), t.Day(), t.Location())
}
func NewDay(year int, month time.Month, day int, location *time.Location) Day {
	return Day(time.Date(year, month, day, 0, 0, 0, 0, location))
}
func (d Day) AsTime() time.Time {
	return time.Time(d)
}

type DayTotal struct {
	Duration Duration
	Day      Day
	Tasks    []TaskTotal // todo replace with aggregating map
}

func NewDayTotal(t time.Time) DayTotal {
	return DayTotal{Day: GetDay(t)}
}

type TaskTotal struct {
	Duration    Duration
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

func GetTaskTotal(task log.Entry, endTime time.Time) TaskTotal {
	return TaskTotal{
		GetDuration(task.Time, endTime),
		task.TaskName,
		task.TaskId,
		false,
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
		if entry.TaskName == "eod" {
			task = GetTaskTotal(entry, entry.Time)
		} else {
			task = GetTaskTotal(entry, endTime)
		}

		dayTotal.Duration = dayTotal.Duration.Add(task.Duration)
		dayTotal.Tasks = append(dayTotal.Tasks, task)
		if GetDay(entry.Time) != GetDay(endTime) {
			total = append(total, dayTotal)
			dayTotal = NewDayTotal(endTime)
		}
	}
	entry := entries[last(entries)]
	if entry.TaskName == "eod" {
		task = GetTaskTotal(entry, entry.Time)
	} else {
		task = TaskTotal{GetDuration(entry.Time, now), entry.TaskName, entry.TaskId, true}
	}
	dayTotal.Duration = dayTotal.Duration.Add(task.Duration)
	dayTotal.Tasks = append(dayTotal.Tasks, task)

	total = append(total, dayTotal)
	return total
}
