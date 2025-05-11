package actors_test

import (
	"log"
	"testing"
	"time"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSimpleWiresModuleActor_FourWiresMoreThanOneRedOddSerial(t *testing.T) {
	// Arrange
	bomb := entities.NewBomb(valueobject.NewDefaultBombConfig())
	bomb.SerialNumber = "1111"
	simpleWiresModule := entities.NewSimpleWiresModule(bomb)
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
			},
			{
				WireColor: valueobject.SimpleWireColors[1],
			},
			{
				WireColor: valueobject.SimpleWireColors[2],
			},
			{
				WireColor: valueobject.SimpleWireColors[0],
			},
		},
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
			strike:    true,
		},
		{
			desc:      "Cut the second wire",
			wireIndex: 1,
			solved:    false,
			strike:    true,
		},
		{
			desc:      "Cut the last Red wire",
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
				assert.True(t, resp.IsSuccess(), "Step %d: expected success response for correct wire", i+1)
			}

			if resp.IsSuccess() {
				successResp, ok := resp.(actors.SuccessResponse)
				log.Printf("Response: %+v", successResp.Data)
				assert.True(t, ok, "Expected SuccessResponse type")

				result, ok := successResp.Data.(*command.SimpleWiresInputCommandResult)
				assert.True(t, ok, "Expected SimpleWiresInputCommandResult type")

				assert.Equal(t, action.solved, result.Solved, "Step %d: solved state mismatch", i+1)
				assert.Equal(t, action.strike, result.Strike, "Step %d: strike state mismatch", i+1)
			} else {
				errorResp := resp.(actors.ErrorResponse)
				log.Printf("Error response: %+v", errorResp)
			}
		})
	}

	// Verify final state
	assert.True(t, specifiedModule.GetModuleState().MarkSolved, "Module should be solved at the end of the test")
	t.Logf("Final state: %s", specifiedModule)
}
