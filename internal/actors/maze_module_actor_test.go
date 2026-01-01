package actors_test

import (
	"log"
	"testing"
	"time"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMazeModuleActor_NavigationNoWall(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	mazeModule := entities.NewMazeModule(rng)
	mazeModule.SetBomb(bomb)
	mazeModuleActor := actors.NewMazeModuleActor(mazeModule)
	mazeModuleActor.Start() // Start the actor to process messages
	defer mazeModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := uuid.New()

	actions := []struct {
		desc      string
		direction valueobject.CardinalDirection
		strike    bool
		solved    bool
	}{
		{
			desc:      "Move left",
			direction: valueobject.West,
		},
		{
			desc:      "Move down",
			direction: valueobject.South,
		},
		{
			desc:      "Move down",
			direction: valueobject.South,
		},
		{
			desc:      "Move down and touch goal",
			direction: valueobject.South,
			solved:    true,
		},
		// TODO: Add expectErr for invalid input after solving
		// 	{
		// 		desc:      "Move up (invalid after solving)",
		// 		direction: valueobject.North,
		// 		solved:    true,
		// 		expectErr: true,
		// 	},
	}

	for i, action := range actions {
		t.Run(action.desc, func(t *testing.T) {
			// Act
			cmd := &command.MazeCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				Direction: action.direction,
			}

			respChan := make(chan actors.Response, 1)

			mazeModuleActor.Send(actors.ModuleCommandMessage{
				Command:         cmd,
				ResponseChannel: respChan,
			})

			var resp actors.Response
			select {
			case resp = <-respChan:
			case <-time.After(1 * time.Second):
				t.Fatalf("Step %d: timeout waiting for response", i+1)
			}

			// Assert
			if successResp, ok := resp.(actors.SuccessResponse); ok {
				result, ok := successResp.Data.(*command.MazeInputCommandResult)
				if !ok {
					t.Fatalf("Expected MazeInputCommandResult type")
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

	log.Printf("Final state: %v", mazeModuleActor.GetModule())
}

func TestMazeModuleActor_NavigationHitWall(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test-walls")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	mazeModule := entities.NewMazeModule(rng)
	mazeModule.SetBomb(bomb)
	mazeModuleActor := actors.NewMazeModuleActor(mazeModule)
	mazeModuleActor.Start() // Start the actor to process messages
	defer mazeModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := uuid.New()

	actions := []struct {
		desc      string
		direction valueobject.CardinalDirection
		strike    bool
		solved    bool
	}{
		{
			desc:      "Move right",
			direction: valueobject.East,
			strike:    true,
		},
		{
			desc:      "Move up",
			direction: valueobject.North,
			strike:    true,
		},
		{
			desc:      "Move down",
			direction: valueobject.South,
		},
		{
			desc:      "Move right",
			direction: valueobject.East,
			solved:    false,
		},
	}

	for i, action := range actions {
		t.Run(action.desc, func(t *testing.T) {
			// Act
			cmd := &command.MazeCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				Direction: action.direction,
			}

			respChan := make(chan actors.Response, 1)

			mazeModuleActor.Send(actors.ModuleCommandMessage{
				Command:         cmd,
				ResponseChannel: respChan,
			})

			var resp actors.Response
			select {
			case resp = <-respChan:
			case <-time.After(1 * time.Second):
				t.Fatalf("Step %d: timeout waiting for response", i+1)
			}

			// Assert
			if successResp, ok := resp.(actors.SuccessResponse); ok {
				result, ok := successResp.Data.(*command.MazeInputCommandResult)
				if !ok {
					t.Fatalf("Expected MazeInputCommandResult type")
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

	log.Printf("Final state: %v", mazeModuleActor.GetModule())
}
