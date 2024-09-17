package log

import (
	"fmt"
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
	return fmt.Sprintf("Entry(%s,%s,%s)", entry.Time.String(), entry.TaskName, entry.TaskId)
}

func Parse(argv []string) Entry {
	input := strings.Join(argv, " ")
	entry := NewLogEntry(input).
		parseDate().
		parseFrontDate().
		ParseTime().
		parseTaskId()
	return entry
}

