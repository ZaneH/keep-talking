package command

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type ModuleInputCommand interface {
	GetSessionID() uuid.UUID
	GetModulePosition() valueobject.ModulePosition
}

type BaseModuleInputCommand struct {
	SessionID      uuid.UUID
	ModulePosition valueobject.ModulePosition
}

func (c *BaseModuleInputCommand) GetSessionID() uuid.UUID {
	return c.SessionID
}

func (c *BaseModuleInputCommand) GetModulePosition() valueobject.ModulePosition {
	return c.ModulePosition
}
