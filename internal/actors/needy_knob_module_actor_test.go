package actors_test

import (
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

func TestNeedyKnobModuleActor_RotateDialFromNorth(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("knob_test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	// Set initial dial direction to North
	knobModule.SetState(entities.NeedyKnobState{
		DisplayedPattern: [][]bool{
			{false, false, true, false, true, true},
			{true, true, true, true, false, true},
		},
		DialDirection:      valueobject.North,
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	cmd := &command.NeedyKnobCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	knobModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	assert.True(t, resp.IsSuccess(), "Expected success response")

	if successResp, ok := resp.(actors.SuccessResponse); ok {
		result, ok := successResp.Data.(*command.NeedyKnobCommandResult)
		assert.True(t, ok, "Expected NeedyKnobCommandResult type")
		assert.Equal(t, valueobject.East, result.DialDirection, "Dial should rotate from North to East")
	}
}

func TestNeedyKnobModuleActor_RotateDialFromEast(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("knob_east")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	knobModule.SetState(entities.NeedyKnobState{
		DisplayedPattern: [][]bool{
			{true, false, true, true, true, true},
			{true, true, true, false, true, false},
		},
		DialDirection:      valueobject.East,
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	cmd := &command.NeedyKnobCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	knobModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	assert.True(t, resp.IsSuccess(), "Expected success response")

	if successResp, ok := resp.(actors.SuccessResponse); ok {
		result, ok := successResp.Data.(*command.NeedyKnobCommandResult)
		assert.True(t, ok, "Expected NeedyKnobCommandResult type")
		assert.Equal(t, valueobject.South, result.DialDirection, "Dial should rotate from East to South")
	}
}

func TestNeedyKnobModuleActor_RotateDialFromSouth(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("knob_south")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	knobModule.SetState(entities.NeedyKnobState{
		DisplayedPattern: [][]bool{
			{false, true, true, false, false, true},
			{true, true, true, true, false, true},
		},
		DialDirection:      valueobject.South,
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	cmd := &command.NeedyKnobCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	knobModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	assert.True(t, resp.IsSuccess(), "Expected success response")

	if successResp, ok := resp.(actors.SuccessResponse); ok {
		result, ok := successResp.Data.(*command.NeedyKnobCommandResult)
		assert.True(t, ok, "Expected NeedyKnobCommandResult type")
		assert.Equal(t, valueobject.West, result.DialDirection, "Dial should rotate from South to West")
	}
}

func TestNeedyKnobModuleActor_RotateDialFromWest(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("knob_west")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	knobModule.SetState(entities.NeedyKnobState{
		DisplayedPattern: [][]bool{
			{false, false, false, false, true, false},
			{true, false, false, true, true, true},
		},
		DialDirection:      valueobject.West,
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	cmd := &command.NeedyKnobCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	knobModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	assert.True(t, resp.IsSuccess(), "Expected success response")

	if successResp, ok := resp.(actors.SuccessResponse); ok {
		result, ok := successResp.Data.(*command.NeedyKnobCommandResult)
		assert.True(t, ok, "Expected NeedyKnobCommandResult type")
		assert.Equal(t, valueobject.North, result.DialDirection, "Dial should rotate from West back to North")
	}
}

func TestNeedyKnobModuleActor_FullRotationCycle(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("knob_cycle")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	knobModule.SetState(entities.NeedyKnobState{
		DisplayedPattern: [][]bool{
			{false, false, true, false, true, true},
			{true, true, true, true, false, true},
		},
		DialDirection:      valueobject.North,
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	expectedDirections := []valueobject.CardinalDirection{
		valueobject.East,  // After 1st rotation
		valueobject.South, // After 2nd rotation
		valueobject.West,  // After 3rd rotation
		valueobject.North, // After 4th rotation (back to start)
	}

	for i, expected := range expectedDirections {
		cmd := &command.NeedyKnobCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
		}

		respChan := make(chan actors.Response, 1)
		knobModuleActor.Send(actors.ModuleCommandMessage{
			Command:         cmd,
			ResponseChannel: respChan,
		})

		var resp actors.Response
		select {
		case resp = <-respChan:
		case <-time.After(1 * time.Second):
			t.Fatalf("Timeout waiting for response on rotation %d", i+1)
		}

		assert.True(t, resp.IsSuccess(), "Expected success response on rotation %d", i+1)

		if successResp, ok := resp.(actors.SuccessResponse); ok {
			result, ok := successResp.Data.(*command.NeedyKnobCommandResult)
			assert.True(t, ok, "Expected NeedyKnobCommandResult type")
			assert.Equal(t, expected, result.DialDirection, "Rotation %d: dial should be %v", i+1, expected)
		}
	}
}

func TestNeedyKnobModuleActor_NoStrikeOnRotation(t *testing.T) {
	// Arrange - rotating the dial never causes a strike
	rng := services.NewSeededRNGFromString("knob_nostrike")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	// Rotate multiple times - should never strike
	for i := 0; i < 5; i++ {
		cmd := &command.NeedyKnobCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
		}

		respChan := make(chan actors.Response, 1)
		knobModuleActor.Send(actors.ModuleCommandMessage{
			Command:         cmd,
			ResponseChannel: respChan,
		})

		var resp actors.Response
		select {
		case resp = <-respChan:
		case <-time.After(1 * time.Second):
			t.Fatalf("Timeout waiting for response on rotation %d", i+1)
		}

		if successResp, ok := resp.(actors.SuccessResponse); ok {
			result, ok := successResp.Data.(*command.NeedyKnobCommandResult)
			assert.True(t, ok, "Expected NeedyKnobCommandResult type")
			assert.False(t, result.Strike, "Rotation should never cause a strike")
		}
	}
}

func TestNeedyKnobModuleActor_NeverSolvesAsNeedy(t *testing.T) {
	// Arrange - Needy modules never "solve"
	rng := services.NewSeededRNGFromString("knob_needy")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	// Rotate multiple times - should never be marked as solved
	for i := 0; i < 10; i++ {
		cmd := &command.NeedyKnobCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
		}

		respChan := make(chan actors.Response, 1)
		knobModuleActor.Send(actors.ModuleCommandMessage{
			Command:         cmd,
			ResponseChannel: respChan,
		})

		var resp actors.Response
		select {
		case resp = <-respChan:
		case <-time.After(1 * time.Second):
			t.Fatalf("Timeout waiting for response on rotation %d", i+1)
		}

		if successResp, ok := resp.(actors.SuccessResponse); ok {
			result, ok := successResp.Data.(*command.NeedyKnobCommandResult)
			assert.True(t, ok, "Expected NeedyKnobCommandResult type")
			assert.False(t, result.Solved, "Needy module should never be marked as solved")
		}
	}
}

func TestNeedyKnobModuleActor_DisplayedPatternPreserved(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("knob_pattern")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	expectedPattern := [][]bool{
		{true, false, true, true, true, true},
		{true, true, true, false, true, false},
	}

	knobModule.SetState(entities.NeedyKnobState{
		DisplayedPattern:   expectedPattern,
		DialDirection:      valueobject.North,
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	cmd := &command.NeedyKnobCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	knobModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	assert.True(t, resp.IsSuccess(), "Expected success response")

	if successResp, ok := resp.(actors.SuccessResponse); ok {
		result, ok := successResp.Data.(*command.NeedyKnobCommandResult)
		assert.True(t, ok, "Expected NeedyKnobCommandResult type")
		assert.Equal(t, expectedPattern, result.DisplayedPattern, "Displayed pattern should be preserved in response")
	}
}

func TestNeedyKnobModuleActor_InvalidCommandType(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	knobModule := entities.NewNeedyKnobModule(rng)
	knobModule.SetBomb(bomb)

	knobModuleActor := actors.NewNeedyKnobModuleActor(knobModule)
	knobModuleActor.Start()
	defer knobModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := knobModule.GetModuleID()

	// Send wrong command type
	cmd := &command.PasswordSubmitCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	knobModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	assert.False(t, resp.IsSuccess(), "Expected error for invalid command type")
}
