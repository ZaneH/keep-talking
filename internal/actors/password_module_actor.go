package actors

import (
	"context"
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type PasswordModuleActor struct {
	module *entities.PasswordModule
}

func NewPasswordModuleActor(solution string) *PasswordModuleActor {
	module := entities.NewPasswordModule(solution)
	return &PasswordModuleActor{
		module: module,
	}
}

func (a *PasswordModuleActor) GetModule() common.Module {
	return a.module
}

func (a *PasswordModuleActor) GetModuleID() uuid.UUID {
	return a.module.ModuleID
}

func (a *PasswordModuleActor) ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error) {
	switch c := cmd.(type) {
	case *command.PasswordLetterChangeCommand:
		if c.Direction == valueobject.Increment {
			a.module.IncrementLetterOption(c.LetterIndex)
		}
	case *command.PasswordSubmitCommand:
		a.module.CheckPassword()
	default:
		return nil, errors.New("unsupported command for password module")
	}

	return nil, nil
}
