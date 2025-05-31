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

func TestPasswordModuleActor_LetterChange(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	passwordModule := entities.NewPasswordModule(rng, nil)
	passwordModule.SetBomb(bomb)
	passwordModuleActor := actors.NewPasswordModuleActor(passwordModule)
	passwordModuleActor.Start() // Start the actor to process messages
	defer passwordModuleActor.Stop()

	var specifiedModule *entities.PasswordModule
	if module, ok := passwordModuleActor.GetModule().(*entities.PasswordModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to PasswordModule")
	}

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := uuid.New()

	actions := []struct {
		desc      string
		index     int
		direction valueobject.IncrementDecrement
		expectErr bool
	}{
		{
			desc:      "Increment first letter",
			index:     0,
			direction: valueobject.Increment,
			expectErr: false,
		},
		{
			desc:      "Decrement first letter",
			index:     0,
			direction: valueobject.Decrement,
			expectErr: false,
		},
		{
			desc:      "Increment last letter",
			index:     4,
			direction: valueobject.Increment,
			expectErr: false,
		},
	}

	for i, action := range actions {
		t.Run(action.desc, func(t *testing.T) {
			// Act
			cmd := &command.PasswordLetterChangeCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				LetterIndex: action.index,
				Direction:   action.direction,
			}

			respChan := make(chan actors.Response, 1)

			passwordModuleActor.Send(actors.ModuleCommandMessage{
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

				result, ok := successResp.Data.(*command.PasswordCommandResult)
				assert.True(t, ok, "Expected PasswordLetterChangeCommandResult type")

				assert.Equal(t, specifiedModule.GetModuleState().IsSolved(), result.Solved, "Solved state should match module state")
				assert.False(t, result.Strike, "No strike should be issued for successful operations")
			}
		})
	}

	log.Printf("Final state: %v", passwordModuleActor.GetModule())
}

func TestPasswordModuleActor_ModuleTypeMismatch(t *testing.T) {
	// This test verifies behavior when there's a module type mismatch
	// The actor code assumes the module is a PasswordModule, but we're testing
	// when that's not the case

	// Create a mock actor with incorrect module type for testing
	// This requires a bit of refactoring in the production code to allow
	// for this test setup, but is important to verify error handling

	// For now, we'll simply note that this test would verify that
	// when the underlying module is not a PasswordModule, appropriate
	// error responses are returned
	t.Skip("Module type mismatch test requires production code modifications")
}
