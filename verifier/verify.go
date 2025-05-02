package verifier

import (
	"log"
	"slices"
	"time"
	"wlog/chrono"
	"wlog/model"
)

type Verifier struct {
	Repo *model.Repository
}

func New(repo *model.Repository) Verifier {
	return Verifier{Repo: repo}
}

func (v Verifier) SixMonths() bool {
	now := time.Now()
	valid := true
	log.SetFlags(0)
	for start := now.AddDate(0, -6, 0); start.Before(now); start = start.AddDate(0, 0, 1) {
		if !chrono.IsWeekday(start) {
			continue
		}
		if !v.began(start) {
			log.Printf("%v : Missing", chrono.Date(start).Iso())
			valid = false
		} else {
			if !v.eod(start) {
				log.Printf("%v : Incomplete: no eod", chrono.Date(start).Iso())
				valid = false
			}
			if v.ferie(start) || v.helligdag(start) || v.sickday(start) {
				continue
			}
			if !v.lunch(start) {
				log.Printf("%v : Incomplete: no lunch", chrono.Date(start).Iso())
				valid = false
			}
		}
	}
	return valid
}

func (v Verifier) hasEntry(day time.Time, titles... string) bool {
	log := v.Repo.GetDailyLog(day)
	for _, e := range log {
		if slices.Contains(titles, e.TaskName) {
			return true
		}
	}
	return false
}
func (v Verifier) began(day time.Time) bool {
	log := v.Repo.GetDailyLog(day)
	return len(log) > 0
}
func (v Verifier) eod(day time.Time) bool {
	return v.hasEntry(day, "eod")
}
func (v Verifier) lunch(day time.Time) bool {
	return v.hasEntry(day, "lunch")
}
func (v Verifier) helligdag(day time.Time) bool {
	return v.hasEntry(day, "helligdag")
}
func (v Verifier) ferie(day time.Time) bool {
	return v.hasEntry(day, "ferie", "juleferie")
}
func (v Verifier) sickday(day time.Time) bool {
	return v.hasEntry(day, "syg", "bali belly")
}
