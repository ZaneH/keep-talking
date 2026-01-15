package actors

import (
	"errors"

	"github.com/ZaneH/defuse.party-go/internal/application/command"
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
)

type NeedyVentGasModuleActor struct {
	BaseModuleActor
}

func NewNeedyVentGasModuleActor(module *entities.NeedyVentGasModule) *NeedyVentGasModuleActor {
	actor := &NeedyVentGasModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *NeedyVentGasModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *NeedyVentGasModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.NeedyVentGasCommand:
		needyVentGasModule, ok := a.module.(*entities.NeedyVentGasModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		strike, err := needyVentGasModule.PressButton(typedCmd.Input)
		result := &command.NeedyVentGasCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: strike,
			},
			DisplayedQuestion:  needyVentGasModule.GetCurrentQuestion(),
			CountdownStartedAt: needyVentGasModule.State.CountdownStartedAt,
			CountdownDuration:  needyVentGasModule.State.CountdownDuration,
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
			Err: errors.New("unsupported command type for needyVentGas module"),
		}
	}
}
