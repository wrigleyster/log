package verifier

import (
	"log"
	"time"
	"wlog/chrono"
	"wlog/model"
)

type Verifier struct {
	Repo model.Repository
}

func New(repo model.Repository) Verifier {
	return Verifier{Repo: repo}
}

func (v Verifier) SixMonths() bool {
	now := time.Now()
	valid := true
	log.SetFlags(0)
	for start := now.AddDate(0, -6, 0); start.Before(now); start = start.AddDate(0, 0, 1) {
		if !v.weekday(start) {
			continue
		}
		if !v.began(start) {
			log.Printf("%v: Missing data.", chrono.Date(start).Iso())
			valid = false
		} else {
			if !v.eod(start) {
				log.Printf("%v: Incomplete data. eod missing.", chrono.Date(start).Iso())
				valid = false
			}
			if !v.lunch(start) {
				log.Printf("%v: Incomplete data. Forgot lunch.", chrono.Date(start).Iso())
				valid = false
			}
		}
	}
	return valid
}

func (v Verifier) weekday(day time.Time) bool {
	return day.Weekday() != time.Saturday && day.Weekday() != time.Sunday
}
func (v Verifier) began(day time.Time) bool {
	log := v.Repo.GetDailyLog(day)
	return len(log) > 0
}
func (v Verifier) eod(day time.Time) bool {
	log := v.Repo.GetDailyLog(day)
	for _, e := range log {
		if e.TaskName == "eod" {
			return true
		}
	}
	return false
}
func (v Verifier) lunch(day time.Time) bool {
	log := v.Repo.GetDailyLog(day)
	for _, e := range log {
		if e.TaskName == "lunch" {
			return true
		}
	}
	return false
}
