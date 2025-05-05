package actors

import (
	"context"
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type SimpleWireModuleActor struct {
	module *entities.SimpleWireModule
}

func NewSimpleWireModuleActor(module *entities.SimpleWireModule) *SimpleWireModuleActor {
	return &SimpleWireModuleActor{
		module: module,
	}
}

func (a *SimpleWireModuleActor) GetModule() common.Module {
	return a.module
}

func (a *SimpleWireModuleActor) GetModuleID() uuid.UUID {
	return a.module.ModuleID
}

func (a *SimpleWireModuleActor) ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error) {
	switch c := cmd.(type) {
	case *command.CutWireCommand:
		return a.module.CutWire(c.WireIndex)
	default:
		return nil, errors.New("unsupported command for wire module")
	}
}
