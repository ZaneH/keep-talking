package entities

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type Module interface {
	GetModuleID() uuid.UUID
	GetModuleState() ModuleState
	GetType() valueobject.ModuleType
	String() string
	GetBomb() *Bomb
	// AddStrike()
}

type ModuleState interface {
	IsSolved() bool
	MarkAsSolved()
}

type BaseModuleState struct {
	MarkSolved bool
}

func (b BaseModuleState) IsSolved() bool {
	return b.MarkSolved
}

func (b *BaseModuleState) MarkAsSolved() {
	b.MarkSolved = true
}
