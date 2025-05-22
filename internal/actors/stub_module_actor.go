package actors

import "github.com/ZaneH/keep-talking/internal/domain/entities"

type StubModuleActor struct {
	BaseModuleActor
}

func NewStubModuleActor(module entities.Module, bufferSize int) *StubModuleActor {
	return &StubModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, bufferSize),
	}
}
