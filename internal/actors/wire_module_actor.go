package actors

import (
	"context"
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type WireModuleActor struct {
	module   *entities.WireModule
	moduleId uuid.UUID
}

func (w *WireModuleActor) ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error) {
	switch c := cmd.(type) {
	case *command.CutWireCommand:
		return w.module.CutWire(c.WireIndex)
	default:
		return nil, errors.New("unsupported command for wire module")
	}
}
