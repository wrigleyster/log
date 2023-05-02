package formatter

import (
	"fmt"
	"strings"
	"time"
	"wlog/log"
	"wlog/manipulation"
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

func formatDuration(d manipulation.Delta) string {
	if d.TaskId() == "" {
		return fmt.Sprintf("%dh %.2dm %s", d.Duration().Hours(), d.Duration().Minutes(), d.TaskName())
	}
	return fmt.Sprintf("%dh %.2dm %s %s", d.Duration().Hours(), d.Duration().Minutes(), d.TaskId(), d.TaskName())
}
func FormatDurations(entries []log.Entry, now time.Time) string {
	ps := manipulation.ToDeltas(entries)
	var lines []string
	for _, p := range ps {
		if p.IsOpenEnded() {
			lines = append(lines, formatDate(p.StartTime()))
			if p.IsEOD() {
				continue
			}
			p.SetEnd(now)
			lines = append(lines, "+"+formatDuration(p))
		} else {
			lines = append(lines, " "+formatDuration(p))
		}
	}
	return strings.Join(lines, "\n")
}
