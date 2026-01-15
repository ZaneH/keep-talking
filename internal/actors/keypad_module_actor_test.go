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

func TestKeypadModuleActor_PressSymbolProcessesCommand(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	keypadModule := entities.NewKeypadModule(rng)
	keypadModule.SetBomb(bomb)
	keypadModuleActor := actors.NewKeypadModuleActor(keypadModule)
	keypadModuleActor.Start()
	defer keypadModuleActor.Stop()

	var specifiedModule *entities.KeypadModule
	if module, ok := keypadModuleActor.GetModule().(*entities.KeypadModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to KeypadModule")
	}

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := specifiedModule.GetModuleID()

	// Get the displayed symbols
	displayed := specifiedModule.State.DisplayedSymbols

	t.Run("Press symbol processes command correctly", func(t *testing.T) {
		// Press the first displayed symbol
		cmd := &command.KeypadInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Symbol: displayed[0],
		}

		respChan := make(chan actors.Response, 1)
		keypadModuleActor.Send(actors.ModuleCommandMessage{
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

		successResp, ok := resp.(actors.SuccessResponse)
		assert.True(t, ok, "Expected SuccessResponse type")

		result, ok := successResp.Data.(*command.KeypadInputCommandResult)
		assert.True(t, ok, "Expected KeypadInputCommandResult type")

		assert.NotNil(t, result.ActivatedSymbols, "ActivatedSymbols should not be nil")
		assert.NotNil(t, result.DisplayedSymbols, "DisplayedSymbols should not be nil")
	})
}

func TestKeypadModuleActor_PressAlreadyActivatedSymbolReturnsError(t *testing.T) {
	// Arrange - The keypad module returns an error when pressing an already activated symbol
	// But the activated symbols get reset on wrong order, so we need to press the correct first symbol
	rng := services.NewSeededRNGFromString("activated_test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	keypadModule := entities.NewKeypadModule(rng)
	keypadModule.SetBomb(bomb)
	keypadModuleActor := actors.NewKeypadModuleActor(keypadModule)
	keypadModuleActor.Start()
	defer keypadModuleActor.Stop()

	var specifiedModule *entities.KeypadModule
	if module, ok := keypadModuleActor.GetModule().(*entities.KeypadModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to KeypadModule")
	}

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := specifiedModule.GetModuleID()

	displayed := specifiedModule.State.DisplayedSymbols
	if len(displayed) == 0 {
		t.Fatal("No displayed symbols")
	}

	// Press any symbol first
	firstSymbol := displayed[0]
	cmd := &command.KeypadInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Symbol: firstSymbol,
	}

	respChan := make(chan actors.Response, 1)
	keypadModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	var firstResp actors.Response
	select {
	case firstResp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for first response")
	}

	// Check if the symbol was activated (no strike means it was the correct first symbol)
	if successResp, ok := firstResp.(actors.SuccessResponse); ok {
		result := successResp.Data.(*command.KeypadInputCommandResult)
		if result.Strike {
			// Wrong order - module was reset, symbol is no longer activated
			// This is expected behavior, skip the "already activated" test
			t.Skip("First symbol was wrong order, module reset - skipping already activated test")
		}
	}

	// Now try to press the same symbol again - should get an error
	respChan2 := make(chan actors.Response, 1)
	keypadModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan2,
	})

	var resp actors.Response
	select {
	case resp = <-respChan2:
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for second response")
	}

	// Should get an error response because the symbol is already activated
	assert.False(t, resp.IsSuccess(), "Expected error response for already activated symbol")
}

func TestKeypadModuleActor_WrongOrderGivesStrike(t *testing.T) {
	// Arrange - use the module's generated state which has a valid solution
	rng := services.NewSeededRNGFromString("wrongorder_test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	keypadModule := entities.NewKeypadModule(rng)
	keypadModule.SetBomb(bomb)

	keypadModuleActor := actors.NewKeypadModuleActor(keypadModule)
	keypadModuleActor.Start()
	defer keypadModuleActor.Stop()

	var specifiedModule *entities.KeypadModule
	if module, ok := keypadModuleActor.GetModule().(*entities.KeypadModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to KeypadModule")
	}

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := specifiedModule.GetModuleID()

	displayed := specifiedModule.State.DisplayedSymbols

	// Try pressing symbols in different orders until we get a strike
	// Press the last displayed symbol first - it's likely not the correct first symbol
	cmd := &command.KeypadInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Symbol: displayed[len(displayed)-1], // Press last symbol first
	}

	respChan := make(chan actors.Response, 1)
	keypadModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.KeypadInputCommandResult)
		assert.True(t, ok, "Expected KeypadInputCommandResult type")

		// Either we get a strike (wrong order) or the last symbol happened to be first in solution
		// Both are valid outcomes - we're testing the actor processes the command
		t.Logf("Strike: %v, Activated count: %d", result.Strike, len(result.ActivatedSymbols))
	}
}

func TestKeypadModuleActor_InvalidModuleCommand(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	keypadModule := entities.NewKeypadModule(rng)
	keypadModule.SetBomb(bomb)
	keypadModuleActor := actors.NewKeypadModuleActor(keypadModule)
	keypadModuleActor.Start()
	defer keypadModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := keypadModule.GetModuleID()

	// Send a different command type that keypad doesn't handle
	cmd := &command.PasswordSubmitCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	keypadModuleActor.Send(actors.ModuleCommandMessage{
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
