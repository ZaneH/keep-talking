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
	module := entities.NewPasswordModule(&solution)
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
		var err error
		if c.Direction == valueobject.Increment {
			a.module.IncrementLetterOption(c.LetterIndex)
		} else if c.Direction == valueobject.Decrement {
			a.module.DecrementLetterOption(c.LetterIndex)
		} else {
			err = errors.New("invalid direction for letter change")
		}

		return &command.PasswordLetterChangeCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.IsSolved(),
				Strike: false,
			},
		}, err
	case *command.PasswordSubmitCommand:
		err := a.module.CheckPassword()

		return &command.PasswordSubmitCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.IsSolved(),
				Strike: err != nil,
			},
		}, err
	default:
		return &command.PasswordSubmitCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.IsSolved(),
				Strike: false,
			},
		}, errors.New("unsupported command for password module")
	}
}
