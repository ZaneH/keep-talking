package actors

import (
	"context"
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type PasswordModuleActor struct {
	module   *entities.PasswordModule
	moduleId uuid.UUID
}

func NewPasswordModuleActor(module *entities.PasswordModule) *PasswordModuleActor {
	return &PasswordModuleActor{
		module:   module,
		moduleId: module.ModuleId,
	}
}

func (w *PasswordModuleActor) ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error) {
	switch c := cmd.(type) {
	case *command.SubmitPasswordCommand:
		return w.module.CheckPassword(c.Password)
	default:
		return nil, errors.New("unsupported command for wire module")
	}
}
