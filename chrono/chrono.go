package chrono

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wrigleyster/gorm/util"
	"github.com/wrigleyster/opt"
)

type Date time.Time

func (date Date) At(h int, m int) time.Time {
	d := time.Time(date)
	return time.Date(d.Year(), d.Month(), d.Day(), h, m, 0, 0, d.Location())
}

func (date Date) On(y int, m int, d int) time.Time {
	t := time.Time(date)
	return time.Date(y, time.Month(m), d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

func (date Date) As(zone *time.Location) time.Time {
	d := time.Time(date)
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), zone)
}

func (date Date) Time() string {
	return fmt.Sprintf("%.2d:%.2d", time.Time(date).Hour(), time.Time(date).Minute())
}
func (d Date) Day() string {
	date := time.Time(d)
	return fmt.Sprintf("%s %d %s", date.Weekday().String(), date.Day(), date.Month())
}
func (d Date) Iso() string {
	date := time.Time(d)
	return date.Format("2006.01.02")
}

type Day time.Time // todo consolidate Day and Date

func GetDay(t time.Time) Day {
	return NewDay(t.Year(), t.Month(), t.Day(), t.Location())
}
func NewDay(year int, month time.Month, day int, location *time.Location) Day {
	return Day(time.Date(year, month, day, 0, 0, 0, 0, location))
}
func (d Day) AsTime() time.Time {
	return time.Time(d)
}
func (d Day) Str() string {
	date := time.Time(d)
	return fmt.Sprintf("%s %d %s", date.Weekday().String(), date.Day(), date.Month())
}

type Duration time.Duration

func DurationOf(h, m int) Duration {
	return Duration(time.Duration(60*h+m) * time.Minute)
}

func GetDuration(a, b time.Time) Duration {
	return Duration(b.Sub(a))
}

func (d Duration) Hours() int {
	duration := time.Duration(d)
	hours := duration.Hours()
	return int(hours)
}
func (d Duration) Minutes() int {
	duration := time.Duration(d)
	minutes := duration.Minutes()
	return int(minutes) % 60
}

func (d Duration) Str() string {
	return fmt.Sprintf("%dh %.2dm", d.Hours(), d.Minutes())
}

func (d Duration) Add(e Duration) Duration {
	return Duration(int64(d) + int64(e))
}

/////////////////

func AbsoluteDate(date time.Time, input string) time.Time {
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
func RelativeDate(date time.Time, input string) time.Time {
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
		finaldate = finaldate.Add(-7 * day)
	}
	return finaldate
}

func ParseDate(dateString string, fallback time.Time) time.Time {
	date := RelativeDate(fallback, dateString)
	return AbsoluteDate(date, dateString)
}
func ParseOptionalDate(dateString string, fallback time.Time) opt.Maybe[time.Time] {
	today := Date(time.Now()).At(0, 0)
	date := ParseDate(dateString, today)
	if today == date {
		return opt.No[time.Time]()
	}
	return opt.Some(date)
}

func ParseTime(input string, fallback time.Time) time.Time {
	startTime := strings.Split(input, ":")
	if 2 != len(startTime) {
		return fallback
	}
	hours, err := strconv.Atoi(startTime[0])
	if err != nil {
		return fallback
	}
	minutes, err := strconv.Atoi(startTime[1])
	if err != nil {
		return fallback
	}
	return Date(fallback).At(hours, minutes)
}
func ParseOptionalTime(timeString string, fallback time.Time) opt.Maybe[time.Time] {
	date := ParseTime(timeString, fallback)
	if date == fallback {
		return opt.No[time.Time]()
	}
	return opt.Some(date)
}
