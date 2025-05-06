package actors

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type SimpleWiresModuleActor struct {
	module *entities.SimpleWiresModule
}

func NewSimpleWireModuleActor(module *entities.SimpleWiresModule) *SimpleWiresModuleActor {
	return &SimpleWiresModuleActor{
		module: module,
	}
}

func (a *SimpleWiresModuleActor) GetModuleID() uuid.UUID {
	return a.module.ModuleID
}

func (a *SimpleWiresModuleActor) GetModule() common.Module {
	return a.module
}

func (a *SimpleWiresModuleActor) ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error) {
	switch cmd := cmd.(type) {
	case *command.SimpleWiresInputCommand:
		err := a.module.CutWire(cmd.WireIndex)
		if err != nil {
			return nil, err
		}

		// TODO: Find a better spot to place this
		// if a.module.IsSolved() {
		// 	a.module.State.IsSolved = true
		// }
	}

	return nil, nil
}
