package command

import (
	"github.com/google/uuid"
)

type ModuleInputCommand interface {
	GetSessionID() uuid.UUID
	GetBombID() uuid.UUID
	GetModuleID() uuid.UUID
}

type BaseModuleInputCommand struct {
	SessionID uuid.UUID
	BombID    uuid.UUID
	ModuleID  uuid.UUID
}

type BaseModuleInputCommandResult struct {
	Solved bool
	Strike bool
}

func (c *BaseModuleInputCommand) GetSessionID() uuid.UUID {
	return c.SessionID
}

func (c *BaseModuleInputCommand) GetModuleID() uuid.UUID {
	return c.SessionID
}

func (c *BaseModuleInputCommand) GetBombID() uuid.UUID {
	return c.BombID
}
