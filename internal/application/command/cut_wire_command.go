package command

import (
	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type CutWireCommand struct {
	SessionId      uuid.UUID
	ModulePosition valueobject.ModulePosition
	WireIndex      int32
}

type CutWireResult struct {
	Result *common.InputResult
}
