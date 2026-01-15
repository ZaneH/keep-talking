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

func TestMemoryModuleActor_Stage1_ScreenNumber1_PressPosition2(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("memory_test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	memoryModule := entities.NewMemoryModule(rng)
	memoryModule.SetBomb(bomb)

	// Set state for stage 1, screen number 1 (press position 2)
	memoryModule.SetState(entities.MemoryState{
		ScreenNumber:     1,
		DisplayedNumbers: []int{3, 2, 1, 4}, // Position 2 (index 1) is 2
		Stage:            1,
	})

	memoryModuleActor := actors.NewMemoryModuleActor(memoryModule)
	memoryModuleActor.Start()
	defer memoryModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := memoryModule.GetModuleID()

	// Stage 1, screen 1: press button in position 2 (index 1)
	cmd := &command.MemoryInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		ButtonIndex: 1, // Position 2 is index 1 (0-indexed)
	}

	respChan := make(chan actors.Response, 1)
	memoryModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MemoryInputCommandResult)
		assert.True(t, ok, "Expected MemoryInputCommandResult type")
		assert.False(t, result.Strike, "Expected no strike for correct answer")
		assert.Equal(t, 2, result.Stage, "Expected stage to advance to 2")
	}
}

func TestMemoryModuleActor_Stage1_ScreenNumber3_PressPosition3(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("memory_test2")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	memoryModule := entities.NewMemoryModule(rng)
	memoryModule.SetBomb(bomb)

	// Set state for stage 1, screen number 3 (press position 3)
	memoryModule.SetState(entities.MemoryState{
		ScreenNumber:     3,
		DisplayedNumbers: []int{4, 1, 2, 3}, // Position 3 (index 2) is 2
		Stage:            1,
	})

	memoryModuleActor := actors.NewMemoryModuleActor(memoryModule)
	memoryModuleActor.Start()
	defer memoryModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := memoryModule.GetModuleID()

	// Stage 1, screen 3: press button in position 3 (index 2)
	cmd := &command.MemoryInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		ButtonIndex: 2, // Position 3 is index 2
	}

	respChan := make(chan actors.Response, 1)
	memoryModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MemoryInputCommandResult)
		assert.True(t, ok, "Expected MemoryInputCommandResult type")
		assert.False(t, result.Strike, "Expected no strike for correct answer")
		assert.Equal(t, 2, result.Stage, "Expected stage to advance to 2")
	}
}

func TestMemoryModuleActor_WrongButtonGivesStrikeAndResetsToStage1(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("memory_wrong")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	memoryModule := entities.NewMemoryModule(rng)
	memoryModule.SetBomb(bomb)

	// Set state for stage 1, screen number 1 (should press position 2)
	memoryModule.SetState(entities.MemoryState{
		ScreenNumber:     1,
		DisplayedNumbers: []int{3, 2, 1, 4},
		Stage:            1,
	})

	memoryModuleActor := actors.NewMemoryModuleActor(memoryModule)
	memoryModuleActor.Start()
	defer memoryModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := memoryModule.GetModuleID()

	// Press wrong position (position 1 instead of 2)
	cmd := &command.MemoryInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		ButtonIndex: 0, // Wrong - should be index 1
	}

	respChan := make(chan actors.Response, 1)
	memoryModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MemoryInputCommandResult)
		assert.True(t, ok, "Expected MemoryInputCommandResult type")
		assert.True(t, result.Strike, "Expected strike for wrong answer")
		assert.Equal(t, 1, result.Stage, "Expected stage to reset to 1")
	}
}

func TestMemoryModuleActor_InvalidButtonIndex(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("memory_invalid")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	memoryModule := entities.NewMemoryModule(rng)
	memoryModule.SetBomb(bomb)

	memoryModuleActor := actors.NewMemoryModuleActor(memoryModule)
	memoryModuleActor.Start()
	defer memoryModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := memoryModule.GetModuleID()

	// Invalid button index (out of range)
	cmd := &command.MemoryInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		ButtonIndex: 5, // Invalid - only 0-3 are valid
	}

	respChan := make(chan actors.Response, 1)
	memoryModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	assert.False(t, resp.IsSuccess(), "Expected error for invalid button index")
}

func TestMemoryModuleActor_Stage2_ScreenNumber1_PressLabel4(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("memory_stage2")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	memoryModule := entities.NewMemoryModule(rng)
	memoryModule.SetBomb(bomb)

	// Set state for stage 2, screen number 1 (press button labeled 4)
	memoryModule.SetState(entities.MemoryState{
		ScreenNumber:     1,
		DisplayedNumbers: []int{3, 4, 1, 2}, // Button 4 is at index 1
		Stage:            2,
	})

	memoryModuleActor := actors.NewMemoryModuleActor(memoryModule)
	memoryModuleActor.Start()
	defer memoryModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := memoryModule.GetModuleID()

	// Stage 2, screen 1: press button labeled 4 (at index 1)
	cmd := &command.MemoryInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		ButtonIndex: 1, // Position of button labeled 4
	}

	respChan := make(chan actors.Response, 1)
	memoryModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MemoryInputCommandResult)
		assert.True(t, ok, "Expected MemoryInputCommandResult type")
		assert.False(t, result.Strike, "Expected no strike for correct answer")
		assert.Equal(t, 3, result.Stage, "Expected stage to advance to 3")
	}
}

func TestMemoryModuleActor_Stage2_ScreenNumber3_PressPosition1(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("memory_stage2_screen3")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	memoryModule := entities.NewMemoryModule(rng)
	memoryModule.SetBomb(bomb)

	// Set state for stage 2, screen number 3 (press position 1)
	memoryModule.SetState(entities.MemoryState{
		ScreenNumber:     3,
		DisplayedNumbers: []int{2, 4, 1, 3},
		Stage:            2,
	})

	memoryModuleActor := actors.NewMemoryModuleActor(memoryModule)
	memoryModuleActor.Start()
	defer memoryModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := memoryModule.GetModuleID()

	// Stage 2, screen 3: press button in position 1 (index 0)
	cmd := &command.MemoryInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		ButtonIndex: 0, // Position 1 is index 0
	}

	respChan := make(chan actors.Response, 1)
	memoryModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MemoryInputCommandResult)
		assert.True(t, ok, "Expected MemoryInputCommandResult type")
		assert.False(t, result.Strike, "Expected no strike for correct answer")
		assert.Equal(t, 3, result.Stage, "Expected stage to advance to 3")
	}
}

func TestMemoryModuleActor_InvalidCommandType(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	memoryModule := entities.NewMemoryModule(rng)
	memoryModule.SetBomb(bomb)

	memoryModuleActor := actors.NewMemoryModuleActor(memoryModule)
	memoryModuleActor.Start()
	defer memoryModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := memoryModule.GetModuleID()

	// Send wrong command type
	cmd := &command.PasswordSubmitCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	memoryModuleActor.Send(actors.ModuleCommandMessage{
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

func TestMemoryModuleActor_NegativeButtonIndex(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("memory_negative")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	memoryModule := entities.NewMemoryModule(rng)
	memoryModule.SetBomb(bomb)

	memoryModuleActor := actors.NewMemoryModuleActor(memoryModule)
	memoryModuleActor.Start()
	defer memoryModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := memoryModule.GetModuleID()

	// Negative button index
	cmd := &command.MemoryInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		ButtonIndex: -1, // Invalid negative index
	}

	respChan := make(chan actors.Response, 1)
	memoryModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for response")
	}

	assert.False(t, resp.IsSuccess(), "Expected error for negative button index")
}
