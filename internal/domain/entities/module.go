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

type ModuleState struct {
	MarkSolved bool
}
