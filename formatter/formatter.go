package formatter

import (
	"fmt"
	"strings"
	"time"
	"wlog/list"
	"wlog/log"
	"wlog/manipulation"
)

type Order int8

const (
	Ascending Order = iota
	Descending
)

func formatDate(date time.Time) string {
	return fmt.Sprintf("%s %d %s", date.Weekday().String(), date.Day(), date.Month())
}
func formatEntry(entry log.Entry) string {
	if entry.TaskId == "" {
		return fmt.Sprintf(" %.2d:%.2d %s", entry.Time.Hour(), entry.Time.Minute(), entry.TaskName)
	}
	return fmt.Sprintf(" %.2d:%.2d %s %s", entry.Time.Hour(), entry.Time.Minute(), entry.TaskId, entry.TaskName)
}

func Format(entries []log.Entry) string {
	var lines []string
	var curDay time.Time
	for i, entry := range entries {
		if i == 0 || !manipulation.IsSameDay(curDay, entry.Time) {
			curDay = entry.Time
			lines = append(lines, formatDate(curDay))
		}

		lines = append(lines, formatEntry(entry))
	}
	return strings.Join(lines, "\n")
}

func formatDayTotal(d manipulation.Duration) string {
	return fmt.Sprintf(", total: %s", d.Str())
}
func formatDay(dayTotal manipulation.DayTotal, order Order) []string {
	var lines []string
	lines = append(lines, formatDate(dayTotal.Day.AsTime())+formatDayTotal(dayTotal.Duration))
	if order == Ascending {
		for _, task := range dayTotal.Tasks {
			if task.IsEOD() {
				continue
			}
			lines = append(lines, task.Str())
		}
	} else {
		for task := range list.InReverse(dayTotal.Tasks) {
			if task.IsEOD() {
				continue
			}
			lines = append(lines, task.Str())
		}
	}
	return lines
}

func formatTotal(total manipulation.Total, order Order) string {
	var lines []string
	if order == Ascending {
		for _, dayTotal := range total {
			lines = append(lines, formatDay(dayTotal, Descending)...)
		}
	} else {
		for dayTotal := range list.InReverse(total) {
			lines = append(lines, formatDay(dayTotal, Descending)...)
		}
	}
	return strings.Join(lines, "\n")
}

func FormatDurations(entries []log.Entry, now time.Time, order Order) string {
	list.Reverse(entries)
	return formatTotal(manipulation.Accumulate(entries, now), order)
}
