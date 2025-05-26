package action

import (
	"fmt"
	"strings"
	"time"
	"wlog/chrono"
	"wlog/log"
	"wlog/model"
)

func UsageAdd(argv Argv) {
	fmt.Println("Add")
	fmt.Printf("\t%s: [SFFEAT] working on x [at 9:30] [yesterday|monday-friday]\n", argv[0])
	fmt.Printf("\t%s: [date] [time] [SFFEAT] working on x\n", argv[0])
	fmt.Printf("\t%s: -a<s|h|f> <time>\t # add sygdom/helligdag/ferie\n", argv[0])
}

func Add(db *model.Repository, argv Argv) {
	msg := log.Parse(argv)
	if time.Now().Sub(msg.Time) < 0 {
		warnOrDie("That event is in the future.")
	}
	println("add")
	log.Add(db, msg)
}
func addSpecialDay(db *model.Repository, argv Argv, dayType string) {
	input := strings.Join(argv, " ") + " " + dayType
	special := log.Parse(strings.Split(input, " "))
	if special.TaskName != dayType {
		die("invalid date")
	}
	special.Time = chrono.Date(special.Time).At(9, 00)
	log.Add(db, special)
	input = strings.Join(argv, " ") + " eod"
	eod := log.Parse(strings.Split(input, " "))
	eod.Time = chrono.Date(eod.Time).At(16, 24)
	log.Add(db, eod)
}
func AddSickDay(db *model.Repository, argv Argv) {
	addSpecialDay(db, argv, "syg")
}
func AddHelligdag(db *model.Repository, argv Argv) {
	addSpecialDay(db, argv, "helligdag")
}
func AddFerie(db *model.Repository, argv Argv) {
	addSpecialDay(db, argv, "ferie")
}
