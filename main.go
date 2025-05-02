package main

import (
	"fmt"
	"github.com/wrigleyster/opt"
	"os"
	"wlog/action"
	"wlog/model"
)

func getDb() model.Repository {
	db := model.Seed("sqlite.db")
	db.CleanChildlessParents()
	return db
}

func printUsage(_ *model.Repository, argv action.Argv) {
	action.UsageAdd(argv)
	action.UsageDelete(argv)
	action.UsageList(argv)
	action.UsageSetId(argv)
	fmt.Printf("%s: -h\n", argv[0])
	os.Exit(1)
}
func parseArgs(argv action.Argv) (func(db *model.Repository, argv action.Argv), action.Argv) {
	if len(argv) == 0 || argv[0] == "-l" {
		return action.ListLog, os.Args
	} else if argv[0] == "-ld" {
		return action.ListLogDiff, os.Args
	} else if argv[0] == "-ll" {
		return action.ListDailyLog, os.Args
	} else if argv[0] == "-lld" {
		return action.ListDailyLogDiff, os.Args
	} else if argv[0] == "-lt" {
		return action.ListTasks, os.Args
	} else if argv[0] == "-s" {
		return action.SetId, os.Args[1:]
	} else if argv[0] == "-dd" {
		return action.DeleteEntry, os.Args[1:]
	} else if argv[0] == "-x" {
		return action.Verify, os.Args
	} else if argv[0] == "-h" {
		return printUsage, os.Args
	} else if argv[0] == os.Args[0] {
		return printUsage, os.Args
	} else if argv[0] == "-af" {
		return action.AddFerie, os.Args[2:]
	} else if argv[0] == "-ah" {
		return action.AddHelligdag, os.Args[2:]
	} else if argv[0] == "-as" {
		return action.AddSickDay, os.Args[2:]
	} else {
		return action.Add, os.Args[1:]
	}
}

func main() {
	action, argv := parseArgs(os.Args[1:])
	db := getDb()
	action(&db, argv)
}
func M[T, U any](n opt.Maybe[T], f func(T) U) opt.Maybe[U] {
	if n.Exists {
		return opt.Some(f(n.Value))
	}
	return opt.No[U]()
}

type Row interface {
	Next() bool
	Scan(dest ...any) error
	Columns() ([]string, error)
}
