package log

import "strings"

func ParseSet(argv []string) (extId, name string) {
	if len(argv) < 3 || "=" != argv[1] {
		return "", ""
	} else {
		return argv[0], strings.Join(argv[2:], " ")
	}
}
