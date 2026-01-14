package actors

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type PasswordModuleActor struct {
	BaseModuleActor
}

func NewPasswordModuleActor(module *entities.PasswordModule) *PasswordModuleActor {
	actor := &PasswordModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *PasswordModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *PasswordModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.PasswordLetterChangeCommand:
		passwordModule, ok := a.module.(*entities.PasswordModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		var err error
		switch typedCmd.Direction {
		case valueobject.Increment:
			passwordModule.IncrementLetterOption(typedCmd.LetterIndex)
		case valueobject.Decrement:
			passwordModule.DecrementLetterOption(typedCmd.LetterIndex)
		default:
			err = errors.New("unsupported letter change option")
		}

		result := &command.PasswordCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: err != nil,
			},
			Letters: passwordModule.GetCurrentGuess(),
		}

		if err != nil {
			msg.ResponseChannel <- ErrorResponse{
				Err: err,
			}
		} else {
			msg.ResponseChannel <- SuccessResponse{
				Data: result,
			}
		}

	case *command.PasswordSubmitCommand:
		passwordModule, ok := a.module.(*entities.PasswordModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		strike := passwordModule.CheckPassword()
		result := &command.PasswordCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: strike,
			},
			Letters: passwordModule.GetCurrentGuess(),
		}

		msg.ResponseChannel <- SuccessResponse{
			Data: result,
		}

	default:
		msg.ResponseChannel <- ErrorResponse{
			Err: errors.New("unsupported command type for password module"),
		}
	}
}
