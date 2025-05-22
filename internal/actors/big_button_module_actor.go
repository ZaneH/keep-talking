package actors

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

type BigButtonActorError error

var (
	ErrBigButtonUnhandledPressType BigButtonActorError = errors.New("unhandled press type")
)

type BigButtonModuleActor struct {
	BaseModuleActor
}

func NewBigButtonModuleActor(module entities.Module) *BigButtonModuleActor {
	actor := &BigButtonModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *BigButtonModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *BigButtonModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.BigButtonInputCommand:
		buttonModule, ok := a.module.(*entities.BigButtonModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrBigButtonUnhandledPressType,
			}
			return
		}

		stripColor, strike, err := buttonModule.PressButton(typedCmd.PressType)
		result := &command.BigButtonInputCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: err != nil,
			},
		}

		if stripColor != nil {
			result.StripColor = stripColor
		}

		if strike {
			msg.ResponseChannel <- SuccessResponse{
				Data: result,
			}
		} else if err != nil {
			msg.ResponseChannel <- ErrorResponse{
				Err: err,
			}
		} else {
			msg.ResponseChannel <- SuccessResponse{
				Data: result,
			}
		}
	default:
		msg.ResponseChannel <- ErrorResponse{
			Err: ErrInvalidModuleCommand,
		}
	}
}
