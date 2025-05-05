package command

import (
	"github.com/ZaneH/keep-talking/internal/application/common"
)

type CutWireCommand struct {
	BaseModuleInputCommand
	WireIndex int
}

type CutWireResult struct {
	Result *common.InputResult
}
