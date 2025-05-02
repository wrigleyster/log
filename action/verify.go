package action

import (
	"fmt"
	"wlog/model"
	"wlog/verifier"
)

func Verify(db *model.Repository, _ Argv) {
	if verifier.SixMonths(db) {
		fmt.Println("Valid")
	}
}
