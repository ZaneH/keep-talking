package actors_test

import (
	"testing"
	"time"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestSimonModuleActor_TestVowelCompleteSequence(t *testing.T) {
	// Arange
	rng := services.NewSeededRNGFromString("test")
	gameSession, _ := actors.NewGameSessionActor(rng, valueobject.NewEasyGameSessionConfig(""))
	defer gameSession.Stop()

	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	bomb.SerialNumber = "AAA"

	module := entities.NewSimonModule(rng, nil)
	module.SetBomb(bomb)
	module.SetState(entities.SimonState{
		DisplaySequence: []valueobject.Color{
			valueobject.Red,
			valueobject.Green,
			valueobject.Yellow,
		},
	})

	bomb.AddModule(module, valueobject.ModulePosition{
		Row: 0, Column: 0, Face: 0,
	})

	gameSession.Start()

	gameSession.Send(actors.AddBombMessage{
		Bomb:            bomb,
		ResponseChannel: make(chan actors.Response, 1),
	})

	sessionID := gameSession.GetSessionID()

	actions := []struct {
		desc   string
		color  valueobject.Color
		solved bool
		strike bool
	}{
		{
			desc:   "Red->Blue",
			color:  valueobject.Blue,
			solved: false,
			strike: false,
		},
		{
			desc:   "Green->Yellow",
			color:  valueobject.Yellow,
			solved: false,
			strike: false,
		},
		{
			desc:   "Incorrect input (not green)",
			color:  valueobject.Red,
			solved: false,
			strike: true,
		},
	}

	for _, action := range actions {
		t.Run(action.desc, func(t *testing.T) {
			cmd := &command.SimonInputCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bomb.ID,
					ModuleID:  module.GetModuleID(),
				},
				Color: action.color,
			}

			respChan := make(chan actors.Response, 1)

			// Act
			gameSession.Send(actors.ModuleCommandMessage{
				Command:         cmd,
				ResponseChannel: respChan,
			})

			var resp actors.Response
			select {
			case resp = <-respChan:
			case <-time.After(1 * time.Second):
				t.Fatal("Timeout waiting for response")
			}

			// Assert
			if successResp, ok := resp.(actors.SuccessResponse); ok {
				result, ok := successResp.Data.(*command.SimonInputCommandResult)
				if !ok {
					t.Fatalf("Expected SimonInputCommandResult type")
				}

				if action.solved {
					assert.True(t, result.Solved, "Expected module to be solved")
				} else {
					assert.False(t, result.Solved, "Expected module not to be solved")
				}

				if action.strike {
					assert.True(t, result.Strike, "Expected strike")
				} else {
					assert.False(t, result.Strike, "Expected no strike")
				}
			} else {
				t.Fatalf("Expected success response, got error: %v", resp)
			}
		})
	}

	t.Logf("Final state: %s", module.String())
}
