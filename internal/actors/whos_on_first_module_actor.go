package actors

import (
	"errors"

	"github.com/ZaneH/defuse.party-go/internal/application/command"
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
)

type WhosOnFirstModuleActor struct {
	BaseModuleActor
}

func NewWhosOnFirstModuleActor(module entities.Module) *WhosOnFirstModuleActor {
	actor := &WhosOnFirstModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *WhosOnFirstModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *WhosOnFirstModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.WhosOnFirstInputCommand:
		whosOnFirstModule, ok := a.module.(*entities.WhosOnFirstModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		strike, err := whosOnFirstModule.PressWord(typedCmd.Word)
		result := &command.WhosOnFirstInputCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: strike,
			},
			ScreenWord:  whosOnFirstModule.State.ScreenWord,
			ButtonWords: whosOnFirstModule.State.ButtonWords,
			Stage:       whosOnFirstModule.State.Stage,
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
			Err: errors.New("unsupported command type for simple wires module"),
		}
	}
}
