package actors

import (
	"context"
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type PasswordModuleActor struct {
	module *entities.PasswordModule
}

func NewPasswordModuleActor(module *entities.PasswordModule) *PasswordModuleActor {
	return &PasswordModuleActor{
		module: module,
	}
}

func (a *PasswordModuleActor) ProcessCommand(ctx context.Context, cmd command.ModuleInputCommand) (interface{}, error) {
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
