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

func TestSimpleWiresModuleActor_SolveBasic(t *testing.T) {
	// Arrange
	simpleWiresModule := entities.NewSimpleWiresModule()
	simpleWiresModuleActor := actors.NewSimpleWiresModuleActor(simpleWiresModule)
	simpleWiresModuleActor.Start() // Start the actor to process messages
	defer simpleWiresModuleActor.Stop()

	var specifiedModule *entities.SimpleWiresModule
	if module, ok := simpleWiresModuleActor.GetModule().(*entities.SimpleWiresModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to SimpleWiresModule")
	}

	testState := entities.SimpleWiresState{
		Wires: []valueobject.SimpleWire{
			{
				WireColor: valueobject.SimpleWireColors[0],
				IsCut:     false,
			},
			{
				WireColor: valueobject.SimpleWireColors[1],
				IsCut:     false,
			},
			{
				WireColor: valueobject.SimpleWireColors[2],
				IsCut:     false,
			},
			{
				WireColor: valueobject.SimpleWireColors[0],
				IsCut:     false,
			},
		},
		SolutionIndices: []int{0, 3},
	}

	specifiedModule.SetState(testState)

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := uuid.New()

	actions := []struct {
		desc      string
		wireIndex int
		solved    bool
		strike    bool
	}{
		{
			desc:      "Cut the first Red wire",
			wireIndex: 0,
			solved:    false,
			strike:    false,
		},
		{
			desc:      "Cut the second wire",
			wireIndex: 1,
			solved:    false,
			strike:    true,
		},
		{
			desc:      "Cut the third wire",
			wireIndex: 2,
			solved:    false,
			strike:    true,
		},
		{
			desc:      "Cut the second Red wire",
			wireIndex: 3,
			solved:    true,
			strike:    false,
		},
	}

	for i, action := range actions {
		t.Run(action.desc, func(t *testing.T) {
			// Act
			cmd := &command.SimpleWiresInputCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				WireIndex: action.wireIndex,
			}

			respChan := make(chan actors.Response, 1)

			simpleWiresModuleActor.Send(actors.ModuleCommandMessage{
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

			if action.strike {
				assert.False(t, resp.IsSuccess(), "Step %d: expected error response for incorrect wire", i+1)
			} else {
				assert.True(t, resp.IsSuccess(), "Step %d: expected success response for correct wire", i+1)
			}

			if resp.IsSuccess() {
				successResp, ok := resp.(actors.SuccessResponse)
				assert.True(t, ok, "Expected SuccessResponse type")

				result, ok := successResp.Data.(*command.SimpleWiresInputCommandResult)
				assert.True(t, ok, "Expected SimpleWiresInputCommandResult type")

				assert.Equal(t, action.solved, result.Solved, "Step %d: solved state mismatch", i+1)
				assert.Equal(t, action.strike, result.Strike, "Step %d: strike state mismatch", i+1)
			}
		})
	}

	// Verify final state
	assert.True(t, specifiedModule.IsSolved(), "Module should be solved at the end of the test")
	t.Logf("Final state: %s", specifiedModule)
}
