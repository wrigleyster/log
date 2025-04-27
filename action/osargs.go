package action

import (
	"github.com/wrigleyster/opt"
	"strconv"
)

type Argv []string

func (argv Argv) getArg(i int, fallback string) string {
	if len(argv) > i {
		return argv[i]
	}
	return fallback
}
func (argv Argv) getIntArg(i, fallback int) int {
	if len(argv) > i {
		if i, e := strconv.Atoi(argv[i]); e == nil {
			return i
		}
	}
	return fallback
}
func (argv Argv) getOptionalIntArg(i, fallback int) opt.Maybe[int] {
	if len(argv) > i {
		if i, e := strconv.Atoi(argv[i]); e == nil {
			return opt.Some(i)
		}
		return opt.No[int]()
	}
	return opt.Some(fallback)
}
