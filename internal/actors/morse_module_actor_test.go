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

func TestMorseModuleActor_ChangeFrequencyIncrement(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("morse_test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	morseModule := entities.NewMorseModule(rng)
	morseModule.SetBomb(bomb)

	initialFreqIdx := morseModule.State.SelectedFrequencyIdx

	morseModuleActor := actors.NewMorseModuleActor(morseModule)
	morseModuleActor.Start()
	defer morseModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := morseModule.GetModuleID()

	cmd := &command.MorseChangeFrequencyCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Direction: valueobject.Increment,
	}

	respChan := make(chan actors.Response, 1)
	morseModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MorseCommandResult)
		assert.True(t, ok, "Expected MorseCommandResult type")
		assert.False(t, result.Strike, "No strike for frequency change")
		assert.False(t, result.Solved, "Module should not be solved yet")

		// Frequency index should have increased (or stayed at max)
		if initialFreqIdx < 15 { // 16 words total, max index is 15
			assert.Equal(t, initialFreqIdx+1, result.SelectedFrequencyIdx, "Frequency index should have increased")
		}
	}
}

func TestMorseModuleActor_ChangeFrequencyDecrement(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("morse_dec")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	morseModule := entities.NewMorseModule(rng)
	morseModule.SetBomb(bomb)

	// Set a starting index that allows decrement
	morseModule.SetState(entities.MorseState{
		SelectedFrequencyIdx: 5,
		DisplayedFrequency:   3.542,
		DisplayedPattern:     "... .... . .-.. .-..",
	})

	morseModuleActor := actors.NewMorseModuleActor(morseModule)
	morseModuleActor.Start()
	defer morseModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := morseModule.GetModuleID()

	cmd := &command.MorseChangeFrequencyCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Direction: valueobject.Decrement,
	}

	respChan := make(chan actors.Response, 1)
	morseModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MorseCommandResult)
		assert.True(t, ok, "Expected MorseCommandResult type")
		assert.False(t, result.Strike, "No strike for frequency change")
		assert.Equal(t, 4, result.SelectedFrequencyIdx, "Frequency index should have decreased")
	}
}

func TestMorseModuleActor_TxCorrectFrequencySolvesModule(t *testing.T) {
	// Arrange - we need to use the module's generated state and find the correct frequency
	rng := services.NewSeededRNGFromString("morse_solve")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	morseModule := entities.NewMorseModule(rng)
	morseModule.SetBomb(bomb)

	morseModuleActor := actors.NewMorseModuleActor(morseModule)
	morseModuleActor.Start()
	defer morseModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := morseModule.GetModuleID()

	// Get the displayed pattern and find which word/frequency it corresponds to
	pattern := morseModule.State.DisplayedPattern
	morseWordToIdx := map[string]int{
		"... .... . .-.. .-..":      0,  // shell -> 3.505
		".... .- .-.. .-.. ...":     1,  // halls -> 3.515
		"... .-.. .. -.-. -.-":      2,  // slick -> 3.522
		"- .-. .. -.-. -.-":         3,  // trick -> 3.532
		"-... --- -..- . ...":       4,  // boxes -> 3.535
		".-.. . .- -.- ...":         5,  // leaks -> 3.542
		"... - .-. --- -... .":      6,  // strobe -> 3.545
		"-... .. ... - .-. ---":     7,  // bistro -> 3.552
		"..-. .-.. .. -.-. -.-":     8,  // flick -> 3.555
		"-... --- -- -... ...":      9,  // bombs -> 3.565
		"-... .-. . .- -.-":         10, // break -> 3.572
		"-... .-. .. -.-. -.-":      11, // brick -> 3.575
		"... - . .- -.-":            12, // steak -> 3.582
		"... - .. -. --.":           13, // sting -> 3.592
		"...- . -.-. - --- .-.":     14, // vector -> 3.595
		"-... . .- - ...":           15, // beats -> 3.600
	}

	targetIdx, found := morseWordToIdx[pattern]
	if !found {
		t.Fatalf("Unknown morse pattern: %s", pattern)
	}

	// Navigate to the correct frequency
	currentIdx := morseModule.State.SelectedFrequencyIdx
	for currentIdx != targetIdx {
		var direction valueobject.IncrementDecrement
		if currentIdx < targetIdx {
			direction = valueobject.Increment
			currentIdx++
		} else {
			direction = valueobject.Decrement
			currentIdx--
		}

		cmd := &command.MorseChangeFrequencyCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Direction: direction,
		}

		respChan := make(chan actors.Response, 1)
		morseModuleActor.Send(actors.ModuleCommandMessage{
			Command:         cmd,
			ResponseChannel: respChan,
		})

		select {
		case <-respChan:
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for frequency change response")
		}
	}

	// Now transmit at the correct frequency
	txCmd := &command.MorseTxCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	morseModuleActor.Send(actors.ModuleCommandMessage{
		Command:         txCmd,
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
		result, ok := successResp.Data.(*command.MorseCommandResult)
		assert.True(t, ok, "Expected MorseCommandResult type")
		assert.False(t, result.Strike, "No strike for correct frequency")
		assert.True(t, result.Solved, "Module should be solved")
	}
}

func TestMorseModuleActor_TxWrongFrequencyGivesStrike(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("morse_strike")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	morseModule := entities.NewMorseModule(rng)
	morseModule.SetBomb(bomb)

	// Set state where selected frequency does NOT match solution
	// Pattern is "shell" but selected frequency is "beats"
	morseModule.SetState(entities.MorseState{
		SelectedFrequencyIdx: 15, // "beats" = 3.600 MHz
		DisplayedFrequency:   3.600,
		DisplayedPattern:     "... .... . .-.. .-..", // "shell" in morse (solution is 3.505)
	})

	morseModuleActor := actors.NewMorseModuleActor(morseModule)
	morseModuleActor.Start()
	defer morseModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := morseModule.GetModuleID()

	cmd := &command.MorseTxCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	morseModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MorseCommandResult)
		assert.True(t, ok, "Expected MorseCommandResult type")
		assert.True(t, result.Strike, "Expected strike for wrong frequency")
		assert.False(t, result.Solved, "Module should not be solved")
	}
}

func TestMorseModuleActor_FrequencyBoundsMin(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("morse_min")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	morseModule := entities.NewMorseModule(rng)
	morseModule.SetBomb(bomb)

	// Set frequency to minimum
	morseModule.SetState(entities.MorseState{
		SelectedFrequencyIdx: 0,
		DisplayedFrequency:   3.505,
		DisplayedPattern:     "... .... . .-.. .-..",
	})

	morseModuleActor := actors.NewMorseModuleActor(morseModule)
	morseModuleActor.Start()
	defer morseModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := morseModule.GetModuleID()

	// Try to decrement below minimum
	cmd := &command.MorseChangeFrequencyCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Direction: valueobject.Decrement,
	}

	respChan := make(chan actors.Response, 1)
	morseModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MorseCommandResult)
		assert.True(t, ok, "Expected MorseCommandResult type")
		assert.Equal(t, 0, result.SelectedFrequencyIdx, "Frequency should stay at minimum")
	}
}

func TestMorseModuleActor_FrequencyBoundsMax(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("morse_max")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	morseModule := entities.NewMorseModule(rng)
	morseModule.SetBomb(bomb)

	// Set frequency to maximum (15 = "beats")
	morseModule.SetState(entities.MorseState{
		SelectedFrequencyIdx: 15,
		DisplayedFrequency:   3.600,
		DisplayedPattern:     "-... . .- - ...",
	})

	morseModuleActor := actors.NewMorseModuleActor(morseModule)
	morseModuleActor.Start()
	defer morseModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := morseModule.GetModuleID()

	// Try to increment above maximum
	cmd := &command.MorseChangeFrequencyCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Direction: valueobject.Increment,
	}

	respChan := make(chan actors.Response, 1)
	morseModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.MorseCommandResult)
		assert.True(t, ok, "Expected MorseCommandResult type")
		assert.Equal(t, 15, result.SelectedFrequencyIdx, "Frequency should stay at maximum")
	}
}

func TestMorseModuleActor_InvalidCommandType(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	morseModule := entities.NewMorseModule(rng)
	morseModule.SetBomb(bomb)

	morseModuleActor := actors.NewMorseModuleActor(morseModule)
	morseModuleActor.Start()
	defer morseModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := morseModule.GetModuleID()

	// Send wrong command type
	cmd := &command.PasswordSubmitCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	morseModuleActor.Send(actors.ModuleCommandMessage{
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
