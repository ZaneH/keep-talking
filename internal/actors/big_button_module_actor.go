package actors

import (
	"github.com/ZaneH/defuse.party-go/internal/application/command"
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
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
				Err: ErrInvalidModuleType,
			}
			return
		}

		stripColor, strike, err := buttonModule.PressButton(typedCmd.PressType, typedCmd.ReleaseTimestamp)
		result := &command.BigButtonInputCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: err != nil,
			},
		}

		if stripColor != nil {
			result.StripColor = stripColor
		}

		// TODO: Remove strike condition, err is being misused
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
