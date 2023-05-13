package formatter

import (
	"fmt"
	"wlog/manipulation"
)

type DurationView manipulation.Total

func (view DurationView) Format(order Order) string {
	return format(view, order)
}
func (view DurationView) Data() manipulation.Total {
	return manipulation.Total(view)
}
func (_ DurationView) FormatDate(dayTotal manipulation.DayTotal) string {
	return fmt.Sprintf("%s, total: %s", dayTotal.Day.Str(), dayTotal.Duration.Str())
}

func (_ DurationView) AddTask(lines []string, task manipulation.TaskTotal) []string {
	if task.IsEOD() {
		return lines
	}
	openness := ternary(task.IsOpen, "+", " ")
	return append(lines, fmt.Sprintf("%s%s %s", openness, task.Duration.Str(), task.Str()))
}
func ternary[C any](predicate bool, t, f C) C {
	if predicate {
		return t
	}
	return f
}
