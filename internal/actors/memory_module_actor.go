package actors

import (
	"github.com/ZaneH/defuse.party-go/internal/application/command"
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
)

type MemoryModuleActor struct {
	BaseModuleActor
}

func NewMemoryModuleActor(module entities.Module) *MemoryModuleActor {
	actor := &MemoryModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *MemoryModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *MemoryModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch typedCmd := cmd.(type) {
	case *command.MemoryInputCommand:
		memoryModule, ok := a.module.(*entities.MemoryModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		strike, err := memoryModule.PressButton(int(typedCmd.ButtonIndex))
		result := &command.MemoryInputCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: strike,
			},
			DisplayedNumbers: memoryModule.State.DisplayedNumbers,
			ScreenNumber:     memoryModule.State.ScreenNumber,
			Stage:            memoryModule.State.Stage,
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
