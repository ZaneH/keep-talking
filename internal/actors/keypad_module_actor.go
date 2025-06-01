package actors

import (
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

type KeypadModuleActor struct {
	BaseModuleActor
}

func NewKeypadModuleActor(module entities.Module) *KeypadModuleActor {
	actor := &KeypadModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *KeypadModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *KeypadModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.KeypadInputCommand:
		keypadModule, ok := a.module.(*entities.KeypadModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		active, strike, err := keypadModule.PressSymbol(typedCmd.Symbol)
		result := &command.KeypadInputCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: strike,
			},
			DisplayedSymbols: keypadModule.State.DisplayedSymbols,
			ActivatedSymbols: active,
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
