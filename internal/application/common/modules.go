package common

import "github.com/google/uuid"

type Solvable interface {
	Module
	IsSolved() bool
}

type Module interface {
	GetModuleID() uuid.UUID
}
