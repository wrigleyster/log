package log

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"wlog/chrono"
	"wlog/list"
)

type Entry struct {
	Time     time.Time
	TaskName string
	TaskId   string
}

func NewLogEntry(input string) Entry {
	return Entry{time.Now(), input, ""}
}

func (entry Entry) parseTime() Entry {
	now := entry.Time
	words := strings.Split(entry.TaskName, " ")
	if len(words) > 2 && words[len(words)-2] == "at" {
		startTime := strings.Split(list.El(words, -1), ":")
		if 2 != len(startTime) {
			return entry
		}
		hours, err := strconv.Atoi(startTime[0])
		if err != nil {
			return entry
		}
		minutes, err := strconv.Atoi(startTime[1])
		if err != nil {
			return entry
		}
		entry.Time = chrono.Date(now).At(hours, minutes)
		entry.TaskName = strings.Join(list.Sl(words, 0, -2), " ")
	} else {
		entry.Time = now
	}
	return entry

}
func (entry Entry) parseTaskId() Entry {
	words := strings.Split(entry.TaskName, " ")
	for i, v := range words {
		if strings.HasPrefix(v, "SFFEAT") ||
			strings.HasPrefix(v, "SFSTRY") {
			entry.TaskId = v
			entry.TaskName = strings.Join(append(words[:i], words[i+1:]...), " ")
			break
		}

	}
	return entry
}

func (entry Entry) parseDate() Entry {
	words := strings.Split(entry.TaskName, " ")
	newDate := relativeDate(entry.Time, list.El(words, -1))
	if newDate == entry.Time {
		return entry
	}
	entry.TaskName = strings.Join(list.Sl(words, 0, -1), " ")
	entry.Time = newDate
	return entry
}

func (entry Entry) IsEOD() bool {
	return strings.ToLower(entry.TaskName) == "eod"
}

func (entry Entry) Str() string {
	return fmt.Sprintf("Entry(%s,%s,%s)", entry.Time.String(), entry.TaskName, entry.TaskId)
}

func Parse(input string) Entry {
	entry := NewLogEntry(input).
		parseDate().
		parseTime().
		parseTaskId()
	return entry
}

func relativeDate(date time.Time, input string) time.Time {
	tempDate := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), time.UTC)
	day := time.Hour * 24
	switch input {
	case "yesterday":
		tempDate = tempDate.Add(-day)
	case "monday":
		tempDate = tempDate.Truncate(7 * day)
	case "tuesday":
		tempDate = tempDate.Truncate(7 * day).Add(day)
	case "wednesday":
		tempDate = tempDate.Truncate(7 * day).Add(2 * day)
	case "thursday":
		tempDate = tempDate.Truncate(7 * day).Add(3 * day)
	case "friday":
		tempDate = tempDate.Truncate(7 * day).Add(4 * day)
	case "saturday":
		tempDate = tempDate.Truncate(7 * day).Add(5 * day)
	case "sunday":
		tempDate = tempDate.Truncate(7 * day).Add(6 * day)
	default:
		return date
	}
	return time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
}
