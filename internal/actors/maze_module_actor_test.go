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

	var specifiedModule *entities.MazeModule
	if module, ok := mazeModuleActor.GetModule().(*entities.MazeModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to MazeModule")
	}

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := uuid.New()

	actions := []struct {
		desc      string
		direction valueobject.CardinalDirection
		expectErr bool
	}{
		{
			desc:      "Move right",
			direction: valueobject.Right,
			expectErr: false,
		},
		{
			desc:      "Move down",
			direction: valueobject.Down,
			expectErr: false,
		},
		{
			desc:      "Move up",
			direction: valueobject.Up,
			expectErr: false,
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

			// Assert
			var resp actors.Response
			select {
			case resp = <-respChan:
			case <-time.After(1 * time.Second):
				t.Fatalf("Step %d: timeout waiting for response", i+1)
			}

			if action.expectErr {
				assert.False(t, resp.IsSuccess(), "Step %d: expected error response", i+1)
			} else {
				assert.True(t, resp.IsSuccess(), "Step %d: expected success response", i+1)
			}

			if resp.IsSuccess() {
				successResp, ok := resp.(actors.SuccessResponse)
				assert.True(t, ok, "Expected SuccessResponse type")

				result, ok := successResp.Data.(*command.MazeCommandResult)
				assert.True(t, ok, "Expected MazeCommandResult type")

				assert.Equal(t, specifiedModule.GetModuleState().IsSolved(), result.Solved, "Solved state should match module state")
				assert.False(t, result.Strike, "No strike should be issued for successful operations")
			}
		})
	}

	log.Printf("Final state: %v", mazeModuleActor.GetModule())
}
