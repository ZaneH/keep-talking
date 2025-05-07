package entities

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type Module interface {
	GetModuleID() uuid.UUID
	IsSolved() bool
	GetType() valueobject.ModuleType
}

type ModuleState struct {
	MarkSolved bool
}
