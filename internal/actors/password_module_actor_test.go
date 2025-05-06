package actors_test

import (
	"testing"
	"time"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPasswordModuleActor_LetterChange(t *testing.T) {
	// Arrange
	passwordModule := actors.NewPasswordModuleActor()
	passwordModule.Start() // Start the actor to process messages
	defer passwordModule.Stop()

	// We need to cast to get access to the module for test setup
	var specifiedModule *entities.PasswordModule
	if module, ok := passwordModule.GetModule().(*entities.PasswordModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to PasswordModule")
	}

	// Session properties
	sessionID := uuid.New()
	modulePosition := valueobject.ModulePosition{
		Row:    0,
		Column: 0,
		Face:   valueobject.Front,
	}

	// Define test actions
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

	// Execute test actions
	for i, action := range actions {
		t.Run(action.desc, func(t *testing.T) {
			cmd := &command.PasswordLetterChangeCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID:      sessionID,
					ModulePosition: modulePosition,
				},
				LetterIndex: action.index,
				Direction:   action.direction,
			}

			respChan := make(chan actors.Response, 1)

			passwordModule.Send(actors.ModuleCommandMessage{
				Command:         cmd,
				ResponseChannel: respChan,
			})

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

				result, ok := successResp.Data.(*command.SimpleWiresInputCommandResult)
				assert.True(t, ok, "Expected SimpleWiresInputCommandResult type")

				// Verify the result data
				assert.Equal(t, specifiedModule.IsSolved(), result.Solved, "Solved state should match module state")
				assert.False(t, result.Strike, "No strike should be issued for successful operations")
			}
		})
	}

	// Test with unsupported command type
	t.Run("Unsupported command type", func(t *testing.T) {
		// Using a command type that's not supported by this module
		cmd := &command.SimpleWiresInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID:      sessionID,
				ModulePosition: modulePosition,
			},
			WireIndex: 0,
		}

		respChan := make(chan actors.Response, 1)

		passwordModule.Send(actors.ModuleCommandMessage{
			Command:         cmd,
			ResponseChannel: respChan,
		})

		var resp actors.Response
		select {
		case resp = <-respChan:
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for response")
		}

		assert.False(t, resp.IsSuccess(), "Expected error response for unsupported command")
	})
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
