package formatter

import (
	"fmt"
	"wlog/chrono"
	"wlog/manipulation"
)

type AgendaView manipulation.Total

func (view AgendaView) Format(order Order) string {
	return format(view, order)
}

func (view AgendaView) Data() manipulation.Total {
	return manipulation.Total(view)
}

func (_ AgendaView) AddTask(lines []string, t manipulation.TaskTotal) []string {
	return append(lines, fmt.Sprintf(" %s %s", chrono.Date(t.StartedAt).Time(), t.Str()))
}
func (_ AgendaView) FormatDate(dayTotal manipulation.DayTotal) string {
	return dayTotal.Day.Str()
}
