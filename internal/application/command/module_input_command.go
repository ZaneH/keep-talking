package command

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type ModuleInputCommand interface {
	GetSessionId() uuid.UUID
	GetModulePosition() valueobject.ModulePosition
}

type BaseModuleInputCommand struct {
	SessionId      uuid.UUID
	ModulePosition valueobject.ModulePosition
}

func (c *BaseModuleInputCommand) GetSessionId() uuid.UUID {
	return c.SessionId
}

func (c *BaseModuleInputCommand) GetModulePosition() valueobject.ModulePosition {
	return c.ModulePosition
}
