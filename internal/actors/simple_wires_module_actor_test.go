package actors_test

import (
	"context"
	"testing"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/ZaneH/keep-talking/internal/testutils"
	"github.com/google/uuid"
)

func TestSimpleWiresModuleActor_SolveBasic(t *testing.T) {
	// Arrange
	wireModule := testutils.NewTestSimpleWireModule(t)

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

	if specifiedModule, ok := wireModule.GetModule().(*entities.SimpleWiresModule); ok {
		specifiedModule.SetState(testState)
	}

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
		t.Logf("Step %d: %s", i+1, action.desc)
		cmd := &command.SimpleWiresInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: uuid.New(),
				ModulePosition: valueobject.ModulePosition{
					Row:    0,
					Column: 0,
					Face:   valueobject.Front,
				},
			},
			WireIndex: action.wireIndex,
		}

		// Act
		result, err := wireModule.ProcessCommand(context.Background(), cmd)
		if specificModule, ok := wireModule.GetModule().(*entities.SimpleWiresModule); ok {
			t.Logf("State after action: %+v", specificModule.State)
		}

		res := result.(*command.SimpleWiresInputCommandResult)

		// Assert
		if res.Solved != action.solved {
			t.Fatalf("Step %d: expected solved to be %v, got %v", i+1, action.solved, res.Solved)
		}

		if res.Strike != action.strike {
			t.Fatalf("Step %d: expected strike to be %v, got %v", i+1, action.strike, res.Strike)
		}

		if action.strike {
			if err == nil {
				t.Fatalf("Step %d: expected an error, got nil", i+1)
			}
		}
	}
}
