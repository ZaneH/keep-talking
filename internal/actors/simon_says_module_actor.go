package actors

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

type SimonSaysActorError error

var (
	ErrSimonSaysUnhandledPressType SimonSaysActorError = errors.New("unhandled press type")
)

type SimonSaysModuleActor struct {
	BaseModuleActor
}

func NewSimonSaysModuleActor(module entities.Module) *SimonSaysModuleActor {
	actor := &SimonSaysModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *SimonSaysModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *SimonSaysModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.SimonSaysInputCommand:
		simonSaysModule, ok := a.module.(*entities.SimonSaysModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrSimonSaysUnhandledPressType,
			}
			return
		}

		nextSeq, strike, err := simonSaysModule.PressColor(typedCmd.Color)
		result := &command.SimonSaysInputCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: strike,
			},
			NextColorSequence: nextSeq,
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
	default:
		msg.ResponseChannel <- ErrorResponse{
			Err: ErrInvalidModuleCommand,
		}
	}
}
