package actors

import "github.com/ZaneH/defuse.party-go/internal/domain/entities"

type StubModuleActor struct {
	BaseModuleActor
}

func NewStubModuleActor(module entities.Module, bufferSize int) *StubModuleActor {
	return &StubModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, bufferSize),
	}
}
