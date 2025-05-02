package verifier

import (
	"fmt"
	"log"
	"slices"
	"time"
	"wlog/chrono"
	"wlog/model"
)

type Verifier struct {
	date                                                  chrono.Date
	entries                                               []model.LogEntry
	weekday, began, eod, ferie, helligdag, sickday, lunch bool
}

func New(repo *model.Repository, date time.Time) Verifier {
	v := Verifier{date: chrono.Date(date), entries: repo.GetDailyLog(date)}
	v.scan()
	return v
}

func SixMonths(repo *model.Repository) bool {
	now := time.Now()
	valid := true
	log.SetFlags(0)
	for start := now.AddDate(0, -6, 0); start.Before(now); start = start.AddDate(0, 0, 1) {
		v := New(repo, start)
		valid = v.IsValid() && valid
	}
	return valid
}

func is(entry model.LogEntry, titles ...string) bool {
	return slices.Contains(titles, entry.TaskName)
}

func (v *Verifier) IsValid() bool {
	if !v.weekday && !v.began {
		return true
	}
	if !v.began {
		log.Printf("%v : Missing", v.date.Iso())
		return false
	} else {
		if !v.eod {
			log.Printf("%v : Incomplete: no eod", v.date.Iso())
			return false
		}
		if v.ferie || v.helligdag || v.sickday {
			return true
		}
		if !v.lunch {
			log.Printf("%v : Incomplete: no lunch", v.date.Iso())
			return false
		}
	}
	return true
}
func (v *Verifier) scan() {
	v.isWeekday()
	for _, entry := range v.entries {
		v.setBegan()
		v.isEod(entry)
		v.isLunch(entry)
		v.isFerie(entry)
		v.isHelligdag(entry)
		v.isSickday(entry)
	}
}
func (v *Verifier) str() string {
	return fmt.Sprintf("b%v, w%v, h%v, f%v, s%v, e%v, l%v", v.began, v.weekday, v.helligdag, v.ferie, v.sickday, v.eod, v.lunch)
}
func (v *Verifier) setBegan() {
	v.began = true
}
func (v *Verifier) isWeekday() {
	v.weekday = chrono.IsWeekday(time.Time(v.date))
}
func (v *Verifier) isEod(entry model.LogEntry) {
	v.eod = v.eod || is(entry, "eod")
}
func (v *Verifier) isLunch(entry model.LogEntry) {
	v.lunch = v.lunch || is(entry, "lunch")
}
func (v *Verifier) isHelligdag(entry model.LogEntry) {
	v.helligdag = v.helligdag || is(entry, "helligdag")
}
func (v *Verifier) isFerie(entry model.LogEntry) {
	v.ferie = v.ferie || is(entry, "ferie", "juleferie")
}
func (v *Verifier) isSickday(entry model.LogEntry) {
	v.sickday = v.sickday || is(entry, "syg", "bali belly")
}
