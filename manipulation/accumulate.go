package manipulation

import (
	"os"
	"sort"
	"strings"
	"time"
	"wlog/chrono"
	"wlog/log"
)

type Total []DayTotal

type DayTotal struct {
	Duration chrono.Duration
	Day      chrono.Day
	Tasks    []TaskTotal
}

func NewDayTotal(t time.Time) DayTotal {
	return DayTotal{Day: chrono.GetDay(t)}
}

type TaskTotal struct {
	StartedAt   time.Time
	Duration    chrono.Duration
	Name, ExtId string
	IsOpen      bool
}

func (t TaskTotal) IsEOD() bool {
	return strings.ToLower(t.Name) == "eod"
}

func (t TaskTotal) Str() string {
	if t.ExtId == "" {
		return t.Name
	}
	return t.ExtId + " " + t.Name
}

func getTaskTotal(task log.Entry, endTime time.Time, isOpen bool) TaskTotal {
	if task.IsEOD() {
		endTime = task.Time
		isOpen = false
	}
	return TaskTotal{
		task.Time,
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
	total, dayTotal := Total{}, DayTotal{}
	var task TaskTotal
	for i, j := 0, 1; j < len(entries); i, j = j, j+1 {
		entry, endTime := entries[i], entries[j].Time
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
	task = getTaskTotal(entries[last(entries)], now, true)
	dayTotal.Duration = dayTotal.Duration.Add(task.Duration)
	dayTotal.Tasks = append(dayTotal.Tasks, task)

	total = append(total, dayTotal)
	return total
}

type aggregate map[string]TaskTotal

func (a aggregate) asList() []TaskTotal {
	var list []TaskTotal
	for _, v := range a {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].StartedAt.Before(list[j].StartedAt)
	})
	return list
}
func (a aggregate) add(task TaskTotal) {
	if t, present := a[task.Str()]; present {
		t.Duration = t.Duration.Add(task.Duration)
		t.StartedAt = task.StartedAt
        t.IsOpen = t.IsOpen || task.IsOpen
		a[task.Str()] = t
	} else {
		a[task.Str()] = task
	}
}

func Aggregate(entries []log.Entry, now time.Time) Total {
	assertAscending(entries)
	total, dayTotal := Total{}, DayTotal{}
	aggregate := make(aggregate)
	var task TaskTotal
	for i, j := 0, 1; j < len(entries); i, j = j, j+1 {
		entry, endTime := entries[i], entries[j].Time
		if i == 0 {
			dayTotal = NewDayTotal(entry.Time)
		}
		task = getTaskTotal(entry, endTime, false)
		aggregate.add(task)
		dayTotal.Duration = dayTotal.Duration.Add(task.Duration)
		if chrono.GetDay(entry.Time) != chrono.GetDay(endTime) {
			dayTotal.Tasks = aggregate.asList()
			total = append(total, dayTotal)
			dayTotal = NewDayTotal(endTime)
			aggregate = make(map[string]TaskTotal)
		}
	}
	task = getTaskTotal(entries[last(entries)], now, true)
	aggregate.add(task)
	dayTotal.Duration = dayTotal.Duration.Add(task.Duration)
	dayTotal.Tasks = aggregate.asList()
	total = append(total, dayTotal)
	return total
}
