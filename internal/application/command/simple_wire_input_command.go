package command

import (
	"github.com/ZaneH/keep-talking/internal/application/common"
)

type SimpleWireInputCommand struct {
	BaseModuleInputCommand
	WireIndex int
}

type SimpleWireInputResult struct {
	Result *common.InputResult
}
