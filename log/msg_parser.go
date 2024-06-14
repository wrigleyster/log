package log

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"wlog/chrono"
	"wlog/list"

	"github.com/wrigleyster/gorm/util"
)

type Entry struct {
	Time     time.Time
	TaskName string
	TaskId   string
}

func NewLogEntry(input string) Entry {
	return Entry{time.Now(), input, ""}
}
func (entry Entry) parseTime(input string) time.Time {
		startTime := strings.Split(input, ":")
		if 2 != len(startTime) {
			return entry.Time
		}
		hours, err := strconv.Atoi(startTime[0])
		if err != nil {
			return entry.Time
		}
		minutes, err := strconv.Atoi(startTime[1])
		if err != nil {
			return entry.Time
		}
		return chrono.Date(entry.Time).At(hours, minutes)
}
func (entry Entry) ParseTime() Entry {
	now := entry.Time.Truncate(time.Minute)
	words := strings.Split(entry.TaskName, " ")
	if len(words) > 1 && strings.Contains(words[0], ":") && !strings.HasSuffix(words[0], ":"){
		startTime := entry.parseTime(words[0])
		if startTime == entry.Time {
			return entry
		}
		entry.Time = startTime
		entry.TaskName = strings.Join(words[1:], " ")
	} else if len(words) > 2 && words[len(words)-2] == "at" {
		startTime := entry.parseTime(list.El(words, -1))
		if startTime == entry.Time {
			return entry
		}
		entry.Time = startTime
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

func (entry Entry) parseFrontDate() Entry {
	words := strings.Split(entry.TaskName, " ")
	newDate := relativeDate(entry.Time, words[0])
	newDate = absoluteDate(newDate, words[0])
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

func absoluteDate(date time.Time, input string) time.Time {
	numbers := strings.Split(input, ".")
	if len(numbers) < 2 || 3 < len(numbers) {
		return date
	}
	var year, month, day int
	var err error
	if len(numbers) == 2 {
		year = date.Year()
		month, err = strconv.Atoi(numbers[0])
		util.Log(err, "unable to parse month")
		day, err = strconv.Atoi(numbers[1])
		util.Log(err, "unable to parse day")
	} else {
		year, err = strconv.Atoi(numbers[0])
		util.Log(err, "unable to parse year")
		month, err = strconv.Atoi(numbers[1])
		util.Log(err, "unable to parse month")
		day, err = strconv.Atoi(numbers[2])
		util.Log(err, "unable to parse day")
	}
	return time.Date(year, time.Month(month), day, date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())

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

	finaldate := time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), date.Location())
	if finaldate.After(date) {
		finaldate = finaldate.Add(-7*day)
	}
	return finaldate
}
