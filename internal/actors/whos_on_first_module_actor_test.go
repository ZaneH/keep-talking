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

func TestWhosOnFirstModuleActor_CorrectWordAdvancesStage(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("whos_test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	whosOnFirstModule := entities.NewWhosOnFirstModule(rng)
	whosOnFirstModule.SetBomb(bomb)

	// Set a known state for predictable testing
	whosOnFirstModule.SetState(entities.WhosOnFirstState{
		ScreenWord: "YES", // YES -> position 2 (third button, index 2)
		ButtonWords: []string{
			"NOTHING", // 0
			"UHHH",    // 1
			"OKAY",    // 2 - look at this word (index from YES=2)
			"MIDDLE",  // 3
			"LEFT",    // 4
			"PRESS",   // 5
		},
		Stage: 1,
	})

	whosOnFirstModuleActor := actors.NewWhosOnFirstModuleActor(whosOnFirstModule)
	whosOnFirstModuleActor.Start()
	defer whosOnFirstModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := whosOnFirstModule.GetModuleID()

	// For OKAY, the solution priority list starts with:
	// "MIDDLE", "NO", "FIRST", "YES", "UHHH", "NOTHING", "WAIT", "OKAY", "LEFT", "READY", "BLANK", "PRESS", "WHAT", "RIGHT"
	// The earliest word from our buttons that appears in this list is "MIDDLE"
	cmd := &command.WhosOnFirstInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Word: "MIDDLE",
	}

	respChan := make(chan actors.Response, 1)
	whosOnFirstModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.WhosOnFirstInputCommandResult)
		assert.True(t, ok, "Expected WhosOnFirstInputCommandResult type")
		assert.False(t, result.Strike, "Expected no strike for correct answer")
		assert.Equal(t, 2, result.Stage, "Expected stage to advance to 2")
	}
}

func TestWhosOnFirstModuleActor_WrongWordGivesStrikeAndResetsStage(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("whos_wrong")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	whosOnFirstModule := entities.NewWhosOnFirstModule(rng)
	whosOnFirstModule.SetBomb(bomb)

	// Set a known state for predictable testing
	whosOnFirstModule.SetState(entities.WhosOnFirstState{
		ScreenWord: "YES", // YES -> position 2
		ButtonWords: []string{
			"NOTHING", // 0
			"UHHH",    // 1
			"OKAY",    // 2 - look at this word
			"MIDDLE",  // 3
			"LEFT",    // 4
			"PRESS",   // 5
		},
		Stage: 2, // Start at stage 2 to verify reset
	})

	whosOnFirstModuleActor := actors.NewWhosOnFirstModuleActor(whosOnFirstModule)
	whosOnFirstModuleActor.Start()
	defer whosOnFirstModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := whosOnFirstModule.GetModuleID()

	// Press the wrong word (PRESS is later in priority than MIDDLE)
	cmd := &command.WhosOnFirstInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Word: "PRESS", // Wrong - MIDDLE should be pressed first
	}

	respChan := make(chan actors.Response, 1)
	whosOnFirstModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.WhosOnFirstInputCommandResult)
		assert.True(t, ok, "Expected WhosOnFirstInputCommandResult type")
		assert.True(t, result.Strike, "Expected strike for wrong answer")
		assert.Equal(t, 1, result.Stage, "Expected stage to reset to 1")
	}
}

func TestWhosOnFirstModuleActor_CompleteAllThreeStagesSolvesModule(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("whos_solve")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	whosOnFirstModule := entities.NewWhosOnFirstModule(rng)
	whosOnFirstModule.SetBomb(bomb)

	whosOnFirstModuleActor := actors.NewWhosOnFirstModuleActor(whosOnFirstModule)
	whosOnFirstModuleActor.Start()
	defer whosOnFirstModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := whosOnFirstModule.GetModuleID()

	// Complete 3 stages
	for stage := 1; stage <= 3; stage++ {
		module := whosOnFirstModuleActor.GetModule().(*entities.WhosOnFirstModule)

		// Find the correct word to press based on current state
		screenWord := module.State.ScreenWord
		buttonWords := module.State.ButtonWords

		// Get location to look at from screen word (matches entity's whosOnFirstScreenWordsToLocation)
		// Note: Uses Unicode RIGHT SINGLE QUOTATION MARK (U+2019) for apostrophes to match entity
		screenWordToLocation := map[string]int{
			"YES":      2,
			"FIRST":    1,
			"DISPLAY":  5,
			"OKAY":     1,
			"SAYS":     5,
			"NOTHING":  2,
			"":         4,
			"BLANK":    3,
			"NO":       5,
			"LED":      2,
			"LEAD":     5,
			"READ":     3,
			"RED":      3,
			"REED":     4,
			"LEED":     4,
			"HOLD ON":  5,
			"YOU":      4,
			"YOU ARE":  5,
			"YOUR":     3,
			"YOU'RE":   3,
			"UR":       0,
			"THERE":    5,
			"THEY'RE":  4,
			"THEIR":    3,
			"THEY ARE": 2,
			"SEE":      5,
			"C":        1,
			"CEE":      5,
		}

		lookIdx, exists := screenWordToLocation[screenWord]
		if !exists {
			t.Fatalf("Unknown screen word: %s", screenWord)
		}

		lookWord := buttonWords[lookIdx]

		// Get the solution list for the look word
		solutions := map[string][]string{
			"READY":   {"YES", "OKAY", "WHAT", "MIDDLE", "LEFT", "PRESS", "RIGHT", "BLANK", "READY"},
			"FIRST":   {"LEFT", "OKAY", "YES", "MIDDLE", "NO", "RIGHT", "NOTHING", "UHHH", "WAIT", "READY", "BLANK", "WHAT", "PRESS", "FIRST"},
			"NO":      {"BLANK", "UHHH", "WAIT", "FIRST", "WHAT", "READY", "RIGHT", "YES", "NOTHING", "LEFT", "PRESS", "OKAY", "NO", "MIDDLE"},
			"BLANK":   {"WAIT", "RIGHT", "OKAY", "MIDDLE", "BLANK", "PRESS", "READY", "NOTHING", "NO", "WHAT", "LEFT", "UHHH", "YES", "FIRST"},
			"NOTHING": {"UHHH", "RIGHT", "OKAY", "MIDDLE", "YES", "BLANK", "NO", "PRESS", "LEFT", "WHAT", "WAIT", "FIRST", "NOTHING", "READY"},
			"YES":     {"OKAY", "RIGHT", "UHHH", "MIDDLE", "FIRST", "WHAT", "PRESS", "READY", "NOTHING", "YES", "LEFT", "BLANK", "NO", "WAIT"},
			"WHAT":    {"UHHH", "WHAT", "LEFT", "NOTHING", "READY", "BLANK", "MIDDLE", "NO", "OKAY", "FIRST", "WAIT", "YES", "PRESS", "RIGHT"},
			"UHHH":    {"READY", "NOTHING", "LEFT", "WHAT", "OKAY", "YES", "RIGHT", "NO", "PRESS", "BLANK", "UHHH", "MIDDLE", "WAIT", "FIRST"},
			"LEFT":    {"RIGHT", "LEFT", "FIRST", "NO", "MIDDLE", "YES", "BLANK", "WHAT", "UHHH", "WAIT", "PRESS", "READY", "OKAY", "NOTHING"},
			"RIGHT":   {"YES", "NOTHING", "READY", "PRESS", "NO", "WAIT", "WHAT", "RIGHT", "MIDDLE", "LEFT", "UHHH", "BLANK", "OKAY", "FIRST"},
			"MIDDLE":  {"BLANK", "READY", "OKAY", "WHAT", "NOTHING", "PRESS", "NO", "WAIT", "LEFT", "MIDDLE", "RIGHT", "FIRST", "UHHH", "YES"},
			"OKAY":    {"MIDDLE", "NO", "FIRST", "YES", "UHHH", "NOTHING", "WAIT", "OKAY", "LEFT", "READY", "BLANK", "PRESS", "WHAT", "RIGHT"},
			"WAIT":    {"UHHH", "NO", "BLANK", "OKAY", "YES", "LEFT", "FIRST", "PRESS", "WHAT", "WAIT", "NOTHING", "READY", "RIGHT", "MIDDLE"},
			"PRESS":   {"RIGHT", "MIDDLE", "YES", "READY", "PRESS", "OKAY", "NOTHING", "UHHH", "BLANK", "LEFT", "FIRST", "WHAT", "NO", "WAIT"},
			"YOU":     {"SURE", "YOU ARE", "YOUR", "YOU'RE", "NEXT", "UH HUH", "UR", "HOLD", "WHAT?", "YOU", "UH UH", "LIKE", "DONE", "U"},
			"YOU ARE": {"YOUR", "NEXT", "LIKE", "UH HUH", "WHAT?", "DONE", "UH UH", "HOLD", "YOU", "U", "YOU'RE", "SURE", "UR", "YOU ARE"},
			"YOUR":    {"UH UH", "YOU ARE", "UH HUH", "YOUR", "NEXT", "UR", "SURE", "U", "YOU'RE", "YOU", "WHAT?", "HOLD", "LIKE", "DONE"},
			"YOU'RE":  {"YOU", "YOU'RE", "UR", "NEXT", "UH UH", "YOU ARE", "U", "YOUR", "WHAT?", "UH HUH", "SURE", "DONE", "LIKE", "HOLD"},
			"UR":      {"DONE", "U", "UR", "UH HUH", "WHAT?", "SURE", "YOUR", "HOLD", "YOU'RE", "LIKE", "NEXT", "UH UH", "YOU ARE", "YOU"},
			"U":       {"UH HUH", "SURE", "NEXT", "WHAT?", "YOU'RE", "UR", "UH UH", "DONE", "U", "YOU", "LIKE", "HOLD", "YOU ARE", "YOUR"},
			"UH HUH":  {"UH HUH", "YOUR", "YOU ARE", "YOU", "DONE", "HOLD", "UH UH", "NEXT", "SURE", "LIKE", "YOU'RE", "UR", "U", "WHAT?"},
			"UH UH":   {"UR", "U", "YOU ARE", "YOU'RE", "NEXT", "UH UH", "DONE", "YOU", "UH HUH", "LIKE", "YOUR", "SURE", "HOLD", "WHAT?"},
			"WHAT?":   {"YOU", "HOLD", "YOU'RE", "YOUR", "U", "DONE", "UH UH", "LIKE", "YOU ARE", "UH HUH", "UR", "NEXT", "WHAT?", "SURE"},
			"DONE":    {"SURE", "UH HUH", "NEXT", "WHAT?", "YOUR", "UR", "YOU'RE", "HOLD", "LIKE", "YOU", "U", "YOU ARE", "UH UH", "DONE"},
			"NEXT":    {"WHAT?", "UH HUH", "UH UH", "YOUR", "HOLD", "SURE", "NEXT", "LIKE", "DONE", "YOU ARE", "UR", "YOU'RE", "U", "YOU"},
			"HOLD":    {"YOU ARE", "U", "DONE", "UH UH", "YOU", "UR", "SURE", "WHAT?", "YOU'RE", "NEXT", "HOLD", "UH HUH", "YOUR", "LIKE"},
			"SURE":    {"YOU ARE", "DONE", "LIKE", "YOU'RE", "YOU", "HOLD", "UH HUH", "UR", "SURE", "U", "WHAT?", "NEXT", "YOUR", "UH UH"},
			"LIKE":    {"YOU'RE", "NEXT", "U", "UR", "HOLD", "DONE", "UH UH", "WHAT?", "UH HUH", "YOU", "LIKE", "SURE", "YOU ARE", "YOUR"},
		}

		validWords, exists := solutions[lookWord]
		if !exists {
			t.Fatalf("Unknown look word: %s", lookWord)
		}

		// Find the earliest valid word that appears in button words
		var correctWord string
		for _, word := range validWords {
			for _, btn := range buttonWords {
				if word == btn {
					correctWord = word
					break
				}
			}
			if correctWord != "" {
				break
			}
		}

		if correctWord == "" {
			t.Fatalf("Could not find correct word for stage %d", stage)
		}

		cmd := &command.WhosOnFirstInputCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Word: correctWord,
		}

		respChan := make(chan actors.Response, 1)
		whosOnFirstModuleActor.Send(actors.ModuleCommandMessage{
			Command:         cmd,
			ResponseChannel: respChan,
		})

		var resp actors.Response
		select {
		case resp = <-respChan:
		case <-time.After(1 * time.Second):
			t.Fatalf("Timeout waiting for response on stage %d", stage)
		}

		assert.True(t, resp.IsSuccess(), "Expected success response on stage %d", stage)

		if successResp, ok := resp.(actors.SuccessResponse); ok {
			result, ok := successResp.Data.(*command.WhosOnFirstInputCommandResult)
			assert.True(t, ok, "Expected WhosOnFirstInputCommandResult type")
			assert.False(t, result.Strike, "Expected no strike on stage %d", stage)

			if stage == 3 {
				assert.True(t, result.Solved, "Expected module to be solved after stage 3")
			}
		}
	}
}

func TestWhosOnFirstModuleActor_InvalidCommandType(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	whosOnFirstModule := entities.NewWhosOnFirstModule(rng)
	whosOnFirstModule.SetBomb(bomb)

	whosOnFirstModuleActor := actors.NewWhosOnFirstModuleActor(whosOnFirstModule)
	whosOnFirstModuleActor.Start()
	defer whosOnFirstModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := whosOnFirstModule.GetModuleID()

	// Send wrong command type
	cmd := &command.PasswordSubmitCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	whosOnFirstModuleActor.Send(actors.ModuleCommandMessage{
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
