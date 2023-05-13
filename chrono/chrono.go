package chrono

import (
	"fmt"
	"time"
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

func IsSameDay(a time.Time, b time.Time) bool {
	return GetDay(a) == GetDay(b)
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
	if hours > 24 {
		panic("more than 24 hours")
	}
	return int(hours)
}
func (d Duration) Minutes() int {
	duration := time.Duration(d)
	minutes := duration.Minutes()
	if minutes > 24*60 {
		panic("more than 24 hours")
	}
	return int(minutes) % 60
}

func (d Duration) Str() string {
	return fmt.Sprintf("%dh %.2dm", d.Hours(), d.Minutes())
}

func (d Duration) Add(e Duration) Duration {
	return Duration(int64(d) + int64(e))
}
