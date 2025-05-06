package actors_test

import (
	"context"
	"testing"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

func TestPasswordModuleActor_BasicTest(t *testing.T) {
	// Arrange
	actorSystem := actors.NewActorSystem()
	sessionActor, err := actorSystem.CreateGameSession(valueobject.NewDefaultGameConfig())
	if err != nil {
		t.Fatalf("failed to create game session: %v", err)
	}

	passwordModule := actors.NewPasswordModuleActor("three")
	sessionActor.AddModule(passwordModule, valueobject.ModulePosition{
		Row:    0,
		Column: 0,
		Face:   valueobject.Front,
	})

	testCases := []struct {
		desc        string
		letterIndex int
		direction   valueobject.IncrementDecrement
		strike      bool
		solved      bool
	}{
		{
			desc:        "Increment letter moves to the next letter",
			letterIndex: 0,
			direction:   valueobject.Increment,
			strike:      false,
			solved:      false,
		},
		{
			desc:        "Decrement letter moves to the previous letter",
			letterIndex: 0,
			direction:   valueobject.Decrement,
			strike:      false,
			solved:      false,
		},
	}

	for i, actions := range testCases {
		// Act
		result, err := sessionActor.ProcessModuleCommand(context.Background(), &command.PasswordLetterChangeCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID:      sessionActor.GetSessionID(),
				ModulePosition: valueobject.ModulePosition{Row: 0, Column: 0, Face: valueobject.Front},
			},
			LetterIndex: actions.letterIndex,
			Direction:   valueobject.Increment,
		})

		t.Logf("Step %d: %s", i+1, actions.desc)

		// Assert
		if err != nil {
			t.Fatalf("failed to process command: %v", err)
		}

		res := result.(*command.PasswordLetterChangeCommandResult)
		if res.Solved != actions.solved {
			t.Fatalf("expected solved to be %v, got %v", actions.solved, res.Solved)
		}

		if res.Strike != actions.strike {
			t.Fatalf("expected strike to be %v, got %v", actions.strike, res.Strike)
		}
	}

	t.Logf("Final module state: %s", passwordModule.GetModule())
}
