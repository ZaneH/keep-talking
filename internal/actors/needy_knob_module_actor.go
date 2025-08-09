package actors

import (
	"errors"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

type NeedyKnobModuleActor struct {
	BaseModuleActor
}

func NewNeedyKnobModuleActor(module *entities.NeedyKnobModule) *NeedyKnobModuleActor {
	actor := &NeedyKnobModuleActor{
		BaseModuleActor: NewBaseModuleActor(module, 50),
	}

	actor.SetMessageHandler(actor.handleMessage)

	return actor
}

func (a *NeedyKnobModuleActor) handleMessage(msg Message) {
	switch m := msg.(type) {
	case ModuleCommandMessage:
		a.handleModuleCommand(m)
	default:
		a.BaseModuleActor.handleMessage(msg)
	}
}

func (a *NeedyKnobModuleActor) handleModuleCommand(msg ModuleCommandMessage) {
	cmd := msg.Command

	switch cmd.(type) {
	case *command.NeedyKnobCommand:
		needyKnobModule, ok := a.module.(*entities.NeedyKnobModule)
		if !ok {
			msg.GetResponseChannel() <- ErrorResponse{
				Err: ErrInvalidModuleType,
			}
			return
		}

		err := needyKnobModule.RotateDial()
		result := &command.NeedyKnobCommandResult{
			BaseModuleInputCommandResult: command.BaseModuleInputCommandResult{
				Solved: a.module.GetModuleState().IsSolved(),
				Strike: false,
			},
			DisplayedPattern:  needyKnobModule.State.DisplayedPattern,
			DialDirection:     needyKnobModule.State.DialDirection,
			CoundownStartedAt: needyKnobModule.State.CountdownStartedAt,
			CountdownDuration: needyKnobModule.State.CountdownDuration,
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
			Err: errors.New("unsupported command type for needyKnob module"),
		}
	}
}
