package command

import (
	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/google/uuid"
)

type CutWireCommand struct {
	ModuleId  uuid.UUID
	WireIndex uuid.UUID
}

type CutWireResult struct {
	Result *common.InputResult
}
