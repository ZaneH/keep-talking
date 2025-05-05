package actors

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/application/interfaces"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type ModuleActor interface {
	ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error)
	GetModule() interfaces.BombModule
	ModuleId() uuid.UUID
}

type WireModuleActor struct {
	module   *entities.WireModule
	moduleId uuid.UUID
}
