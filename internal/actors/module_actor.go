package actors

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/application/interfaces"
	"github.com/google/uuid"
)

type ModuleActor interface {
	ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error)
	GetModule() interfaces.BombModule
	ModuleId() uuid.UUID
}
