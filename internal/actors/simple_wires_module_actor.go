package actors

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

type SimpleWiresModuleActor struct {
	BaseModuleActor
}

func NewSimpleWiresModuleActor(module entities.Module) *SimpleWiresModuleActor {
	actor := &SimpleWiresModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *SimpleWiresModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *SimpleWiresModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.SimpleWiresInputCommand:
		wiresModule, ok := a.module.(*entities.SimpleWiresModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: errors.New("invalid module type"),
			}
			return
		}

		strike, err := wiresModule.CutWire(typedCmd.WireIndex)
		result := &command.SimpleWiresInputCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().MarkSolved,
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
