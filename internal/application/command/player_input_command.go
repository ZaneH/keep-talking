package command

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type PlayerInputCommand struct {
	SessionId      uuid.UUID
	ModulePosition valueobject.ModulePosition
	Input          interface{}
}
