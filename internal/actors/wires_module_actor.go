package actors

import (
	"errors"

	"github.com/ZaneH/defuse.party-go/internal/application/command"
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
)

type WiresModuleActor struct {
	BaseModuleActor
}

func NewWiresModuleActor(module entities.Module) *WiresModuleActor {
	actor := &WiresModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *WiresModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *WiresModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.WiresInputCommand:
		wiresModule, ok := a.module.(*entities.WiresModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		strike, err := wiresModule.CutWire(typedCmd.WirePosition)
		result := &command.WiresInputCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: strike,
			},
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
