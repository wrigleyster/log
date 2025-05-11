package action

import (
	"wlog/model"
	"wlog/serve"
)

func Serve(_ *model.Repository, _ Argv) {
	serve.Serve()
}
