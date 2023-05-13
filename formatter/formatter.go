package formatter

import (
	"strings"
	"wlog/list"
	"wlog/manipulation"
)

type Order int8

const (
	Ascending Order = iota
	Descending
)

type View interface {
	Data() manipulation.Total
	AddTask(lines []string, task manipulation.TaskTotal) []string
	FormatDate(dayTotal manipulation.DayTotal) string
	Format(order Order) string
}

func format(view View, order Order) string {
	var lines []string
	if order == Ascending {
		for _, dayTotal := range view.Data() {
			lines = formatDay(lines, view, dayTotal, order)
		}
	} else {
		for dayTotal := range list.InReverse(view.Data()) {
			lines = formatDay(lines, view, dayTotal, Descending)
		}
	}
	return strings.Join(lines, "\n")
}
func formatDay(lines []string, view View, dayTotal manipulation.DayTotal, order Order) []string {
	lines = append(lines, view.FormatDate(dayTotal))
	if order == Ascending {
		for _, task := range dayTotal.Tasks {
			lines = view.AddTask(lines, task)
		}
	} else {
		for task := range list.InReverse(dayTotal.Tasks) {
			lines = view.AddTask(lines, task)
		}
	}
	return lines
}
