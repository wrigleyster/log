package manipulation

import (
	"fmt"
	"strings"
	"time"
	"wlog/list"
	"wlog/log"
)

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

type Delta struct {
	a, b log.Entry
}

func ToDeltas(entries []log.Entry) []Delta {
	es := make([]log.Entry, len(entries), len(entries)+1)
	copy(es, entries)
	list.Reverse(es)
	return pair(append(es, log.Entry{}))
}
func pair(entries []log.Entry) []Delta {
	if len(entries) == 1 {
		return make([]Delta, 0, 10)
	} else {
		return append(pair(list.Sl(entries, 1)), Delta{entries[0], entries[1]})
	}
}
func (p *Delta) SetEnd(time time.Time) {
	p.b.Time = time
}
func (p *Delta) IsOpenEnded() bool {
	return p.b == log.Entry{} ||
		!IsSameDay(p.a.Time, p.b.Time)
}
func (p *Delta) IsEOD() bool {
	return strings.ToLower(p.a.TaskName) == "eod"
}
func (p *Delta) StartTime() time.Time {
	return p.a.Time
}
func (p *Delta) Duration() Duration {
	return Duration(p.b.Time.Sub(p.a.Time))
}
func (p *Delta) TaskName() string {
	return p.a.TaskName
}
func (p *Delta) TaskId() string {
	return p.a.TaskId
}

//func aggregateDeltas(){}
