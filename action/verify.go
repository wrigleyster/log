package action

import (
	"fmt"
	"wlog/model"
	"wlog/verifier"
)

func Verify(db *model.Repository, _ Argv) {
	v := verifier.New(db)
	if v.SixMonths() {
		fmt.Println("Valid")
	}
}
