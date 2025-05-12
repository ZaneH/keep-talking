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

type ModuleInputCommandResult interface {
	HasStrike() bool
	IsSolved() bool
}

func (r BaseModuleInputCommandResult) HasStrike() bool {
	return r.Strike
}

func (r BaseModuleInputCommandResult) IsSolved() bool {
	return r.Solved
}

func (c *BaseModuleInputCommand) GetSessionID() uuid.UUID {
	return c.SessionID
}

func (c *BaseModuleInputCommand) GetModuleID() uuid.UUID {
	return c.ModuleID
}

func (c *BaseModuleInputCommand) GetBombID() uuid.UUID {
	return c.BombID
}
