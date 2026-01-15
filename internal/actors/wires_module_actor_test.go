package actors_test

import (
	"testing"
	"time"

	"github.com/ZaneH/defuse.party-go/internal/actors"
	"github.com/ZaneH/defuse.party-go/internal/application/command"
	"github.com/ZaneH/defuse.party-go/internal/domain/entities"
	"github.com/ZaneH/defuse.party-go/internal/domain/services"
	"github.com/ZaneH/defuse.party-go/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWiresModuleActor_FourWiresMoreThanOneRedOddSerial(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	bomb.SerialNumber = "1111"
	simpleWiresModule := entities.NewWiresModule(rng)
	simpleWiresModule.SetBomb(bomb)
	simpleWiresModuleActor := actors.NewWiresModuleActor(simpleWiresModule)
	simpleWiresModuleActor.Start() // Start the actor to process messages
	defer simpleWiresModuleActor.Stop()

	var specifiedModule *entities.WiresModule
	if module, ok := simpleWiresModuleActor.GetModule().(*entities.WiresModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to WiresModule")
	}

	testState := entities.WiresState{
		Wires: []valueobject.Wire{
			{
				// Yellow
				WireColor: valueobject.Yellow,
				Position:  1,
			},
			{
				// Red
				WireColor: valueobject.Red,
				Position:  2,
			},
			{
				// Blue
				WireColor: valueobject.Blue,
				Position:  3,
			},
			{
				// Red
				WireColor: valueobject.Red,
				Position:  4,
			},
		},
	}

	specifiedModule.SetState(testState)

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := uuid.New()

	actions := []struct {
		desc    string
		wirePos int
		solved  bool
		strike  bool
	}{
		{
			desc:    "Cut the first wire (Yellow)",
			wirePos: 1,
			solved:  false,
			strike:  true,
		},
		{
			desc:    "Cut the second wire (Red)",
			wirePos: 2,
			solved:  false,
			strike:  true,
		},
		{
			desc:    "Cut the last Red wire",
			wirePos: 4,
			solved:  true,
			strike:  false,
		},
	}

	for i, action := range actions {
		t.Run(action.desc, func(t *testing.T) {
			// Act
			cmd := &command.WiresInputCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				WirePosition: action.wirePos,
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
				assert.True(t, ok, "Expected SuccessResponse type")

				result, ok := successResp.Data.(*command.WiresInputCommandResult)
				assert.True(t, ok, "Expected WiresInputCommandResult type")

				assert.Equal(t, action.solved, result.Solved, "Step %d: solved state mismatch", i+1)
				assert.Equal(t, action.strike, result.Strike, "Step %d: strike state mismatch", i+1)
			} else {
				t.Errorf("Step %d: expected success response, got error", i+1)
			}
		})
	}

	// Verify final state
	assert.True(t, specifiedModule.GetModuleState().IsSolved(), "Module should be solved at the end of the test")
	t.Logf("Final state: %s", specifiedModule)
}

func TestWiresModuleActor_ThreeWiresNoRed(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	bomb.SerialNumber = "1111"
	wiresModule := entities.NewWiresModule(rng)
	wiresModule.SetBomb(bomb)
	wiresModuleActor := actors.NewWiresModuleActor(wiresModule)
	wiresModuleActor.Start() // Start the actor to process messages
	defer wiresModuleActor.Stop()

	var specifiedModule *entities.WiresModule
	if module, ok := wiresModuleActor.GetModule().(*entities.WiresModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to WiresModule")
	}

	testState := entities.WiresState{
		Wires: []valueobject.Wire{
			{
				// Yellow
				WireColor: valueobject.Yellow,
				Position:  1,
			},
			{
				// Blue
				WireColor: valueobject.Blue,
				Position:  2,
			},
			{
				// Black
				WireColor: valueobject.Black,
				Position:  3,
			},
		},
	}

	specifiedModule.SetState(testState)

	sessionID := uuid.New()
	bombID := bomb.ID
	moduleID := wiresModule.GetModuleID()

	actions := []struct {
		desc    string
		wirePos int
		solved  bool
		strike  bool
	}{
		{
			desc:    "Cut the first wire",
			wirePos: 1,
			solved:  false,
			strike:  true,
		},
		{
			desc:    "Cut the third wire",
			wirePos: 3,
			solved:  false,
			strike:  true,
		},
		{
			desc:    "Cut the second wire",
			wirePos: 2,
			solved:  true,
			strike:  false,
		},
	}

	for i, action := range actions {
		t.Run(action.desc, func(t *testing.T) {
			// Act
			cmd := &command.WiresInputCommand{
				BaseModuleInputCommand: command.BaseModuleInputCommand{
					SessionID: sessionID,
					BombID:    bombID,
					ModuleID:  moduleID,
				},
				WirePosition: action.wirePos,
			}

			respChan := make(chan actors.Response, 1)

			wiresModuleActor.Send(actors.ModuleCommandMessage{
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
			if action.strike {
				assert.True(t, resp.IsSuccess(), "Step %d: expected success response for correct wire", i+1)
			}

			if resp.IsSuccess() {
				successResp, ok := resp.(actors.SuccessResponse)
				assert.True(t, ok, "Expected SuccessResponse type")

				result, ok := successResp.Data.(*command.WiresInputCommandResult)
				assert.True(t, ok, "Expected WiresInputCommandResult type")

				assert.Equal(t, action.solved, result.Solved, "Step %d: solved state mismatch", i+1)
				assert.Equal(t, action.strike, result.Strike, "Step %d: strike state mismatch", i+1)
			} else {
				t.Errorf("Step %d: expected success response, got error", i+1)
			}
		})
	}
}
