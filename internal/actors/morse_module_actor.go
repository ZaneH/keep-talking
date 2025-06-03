package actors

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

type MorseModuleActor struct {
	BaseModuleActor
}

func NewMorseModuleActor(module *entities.MorseModule) *MorseModuleActor {
	actor := &MorseModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *MorseModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *MorseModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.MorseChangeFrequencyCommand:
		morseModule, ok := a.module.(*entities.MorseModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		morseModule.PressChangeFrequency(typedCmd.Direction)

		result := &command.MorseCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: false,
			},
			DisplayedFrequency:   morseModule.GetCurrentFrequency(),
			SelectedFrequencyIdx: morseModule.State.SelectedFrequencyIdx,
		}

		msg.ResponseChannel <- SuccessResponse{
			Data: result,
		}

	case *command.MorseTxCommand:
		morseModule, ok := a.module.(*entities.MorseModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		strike, err := morseModule.PressTx()
		result := &command.MorseCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: strike,
			},
			DisplayedFrequency:   morseModule.GetCurrentFrequency(),
			SelectedFrequencyIdx: morseModule.State.SelectedFrequencyIdx,
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
			Err: errors.New("unsupported command type for password module"),
		}
	}
}
