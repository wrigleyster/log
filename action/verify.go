package action

import (
	"fmt"
	"wlog/model"
	"wlog/verifier"
)

func UsageVerify(argv Argv) {
	fmt.Println("Verify")
	fmt.Printf("\t%s: -cc\n", argv[0])
}

func Verify(db *model.Repository, _ Argv) {
	if verifier.SixMonths(db) {
		fmt.Println("Valid")
	}
}
