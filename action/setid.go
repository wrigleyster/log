package action

import (
	"fmt"
	"os"
	"wlog/log"
	"wlog/model"
)

func UsageSetId(argv Argv) {
	fmt.Println("Set")
	fmt.Printf("\t%s: -s SFFEAT = worked on x [at 9:30] [yesterday|monday-friday]\n", argv[0])
}

func SetId(db *model.Repository, argv Argv) {
	extId, name := log.ParseSet(argv)
	if extId != "" && name != "" {
		tasks := db.TasksByName(name)
		if len(tasks) == 1 {
			tasks[0].ExtId = extId
			db.SaveTask(&tasks[0])
		} else {
			println("Error: multiple tasks named: ", name)
			os.Exit(1)
		}
	} else {
		println("Error: invalid assignment")
		os.Exit(1)
	}
}
