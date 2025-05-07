package entities

import "github.com/ZaneH/keep-talking/internal/domain/valueobject"

type Module interface {
	IsSolved() bool
	GetType() valueobject.ModuleType
}

type ModuleState struct {
	MarkSolved bool
}
