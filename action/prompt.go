package action

import (
	"fmt"
	"github.com/wrigleyster/gorm/util"
	"os"
	"strings"
)

func Prompt(prompt ...string) string {
	for _, msg := range prompt {
		print(msg)
	}
	var reply string
	_, err := fmt.Scanln(&reply)
	util.Log(err)
	return reply
}

func warnOrDie(msg string) {
	response := Prompt("Warning:", msg, "Proceed anyway [y/N]: ")
	response = strings.ToLower(response)
	if response != "y" && response != "yes" {
		os.Exit(1)
	}
}

func die(msg string) {
	println(msg)
	os.Exit(1)
}
