package manipulation

import (
	"fmt"
	"time"
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
