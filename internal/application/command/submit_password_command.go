package command

import (
	"github.com/ZaneH/keep-talking/internal/application/common"
)

type SubmitPasswordCommand struct {
	BaseModuleInputCommand
	Password string
}

type SubmitPasswordResult struct {
	Result *common.InputResult
}
