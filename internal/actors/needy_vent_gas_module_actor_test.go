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

func TestNeedyVentGasModuleActor_VentGasQuestion_YesIsCorrect(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("vent_test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	ventGasModule := entities.NewNeedyVentGasModule(rng)
	ventGasModule.SetBomb(bomb)

	// Set state with "vent gas?" question (answer is yes/true)
	ventGasModule.SetState(entities.NeedyVentGasState{
		DisplayedQuestion:  "vent gas?",
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	ventGasModuleActor := actors.NewNeedyVentGasModuleActor(ventGasModule)
	ventGasModuleActor.Start()
	defer ventGasModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := ventGasModule.GetModuleID()

	// Answer yes (true) to "vent gas?"
	cmd := &command.NeedyVentGasCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Input: true, // Yes
	}

	respChan := make(chan actors.Response, 1)
	ventGasModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.NeedyVentGasCommandResult)
		assert.True(t, ok, "Expected NeedyVentGasCommandResult type")
		assert.False(t, result.Strike, "Expected no strike for correct answer (Yes to vent gas)")
	}
}

func TestNeedyVentGasModuleActor_VentGasQuestion_NoIsIncorrect(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("vent_wrong")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	ventGasModule := entities.NewNeedyVentGasModule(rng)
	ventGasModule.SetBomb(bomb)

	// Set state with "vent gas?" question (answer is yes/true)
	ventGasModule.SetState(entities.NeedyVentGasState{
		DisplayedQuestion:  "vent gas?",
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	ventGasModuleActor := actors.NewNeedyVentGasModuleActor(ventGasModule)
	ventGasModuleActor.Start()
	defer ventGasModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := ventGasModule.GetModuleID()

	// Answer no (false) to "vent gas?" - this is wrong
	cmd := &command.NeedyVentGasCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Input: false, // No - wrong answer
	}

	respChan := make(chan actors.Response, 1)
	ventGasModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.NeedyVentGasCommandResult)
		assert.True(t, ok, "Expected NeedyVentGasCommandResult type")
		assert.True(t, result.Strike, "Expected strike for wrong answer (No to vent gas)")
	}
}

func TestNeedyVentGasModuleActor_CorrectAnswerBasedOnQuestion(t *testing.T) {
	// Arrange - use the module's generated state to determine the correct answer
	rng := services.NewSeededRNGFromString("vent_correct")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	ventGasModule := entities.NewNeedyVentGasModule(rng)
	ventGasModule.SetBomb(bomb)

	ventGasModuleActor := actors.NewNeedyVentGasModuleActor(ventGasModule)
	ventGasModuleActor.Start()
	defer ventGasModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := ventGasModule.GetModuleID()

	// Get the current question and determine the correct answer
	// "vent gas?" -> Yes (true)
	// "detonate?" -> No (false)
	question := ventGasModule.GetCurrentQuestion()
	correctAnswer := question == "vent gas?"

	t.Logf("Question: %s, Correct answer: %v", question, correctAnswer)

	cmd := &command.NeedyVentGasCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Input: correctAnswer,
	}

	respChan := make(chan actors.Response, 1)
	ventGasModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.NeedyVentGasCommandResult)
		assert.True(t, ok, "Expected NeedyVentGasCommandResult type")
		assert.False(t, result.Strike, "Expected no strike for correct answer")
	}
}

func TestNeedyVentGasModuleActor_WrongAnswerGivesStrike(t *testing.T) {
	// Arrange - use the module's generated state and give the wrong answer
	rng := services.NewSeededRNGFromString("vent_wrong")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	ventGasModule := entities.NewNeedyVentGasModule(rng)
	ventGasModule.SetBomb(bomb)

	ventGasModuleActor := actors.NewNeedyVentGasModuleActor(ventGasModule)
	ventGasModuleActor.Start()
	defer ventGasModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := ventGasModule.GetModuleID()

	// Get the current question and give the WRONG answer
	// "vent gas?" -> correct is Yes (true), wrong is No (false)
	// "detonate?" -> correct is No (false), wrong is Yes (true)
	question := ventGasModule.GetCurrentQuestion()
	wrongAnswer := question != "vent gas?" // Opposite of correct answer

	t.Logf("Question: %s, Wrong answer: %v", question, wrongAnswer)

	cmd := &command.NeedyVentGasCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Input: wrongAnswer,
	}

	respChan := make(chan actors.Response, 1)
	ventGasModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.NeedyVentGasCommandResult)
		assert.True(t, ok, "Expected NeedyVentGasCommandResult type")
		assert.True(t, result.Strike, "Expected strike for wrong answer")
	}
}

func TestNeedyVentGasModuleActor_AnswerResetsCountdown(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("vent_countdown")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	ventGasModule := entities.NewNeedyVentGasModule(rng)
	ventGasModule.SetBomb(bomb)

	oldTime := time.Now().Add(-10 * time.Second).Unix()
	ventGasModule.SetState(entities.NeedyVentGasState{
		DisplayedQuestion:  "vent gas?",
		CountdownStartedAt: oldTime,
		CountdownDuration:  30,
	})

	ventGasModuleActor := actors.NewNeedyVentGasModuleActor(ventGasModule)
	ventGasModuleActor.Start()
	defer ventGasModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := ventGasModule.GetModuleID()

	cmd := &command.NeedyVentGasCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		Input: true,
	}

	respChan := make(chan actors.Response, 1)
	ventGasModuleActor.Send(actors.ModuleCommandMessage{
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
		result, ok := successResp.Data.(*command.NeedyVentGasCommandResult)
		assert.True(t, ok, "Expected NeedyVentGasCommandResult type")
		// Countdown should have been reset to a new time
		assert.Greater(t, result.CountdownStartedAt, oldTime, "Countdown should have been reset")
	}
}

func TestNeedyVentGasModuleActor_InvalidCommandType(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	ventGasModule := entities.NewNeedyVentGasModule(rng)
	ventGasModule.SetBomb(bomb)

	ventGasModuleActor := actors.NewNeedyVentGasModuleActor(ventGasModule)
	ventGasModuleActor.Start()
	defer ventGasModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := ventGasModule.GetModuleID()

	// Send wrong command type
	cmd := &command.PasswordSubmitCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
	}

	respChan := make(chan actors.Response, 1)
	ventGasModuleActor.Send(actors.ModuleCommandMessage{
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

func TestNeedyVentGasModuleActor_NeverSolvesAsNeedy(t *testing.T) {
	// Arrange - Needy modules are never "solved" in the traditional sense
	rng := services.NewSeededRNGFromString("vent_needy")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	ventGasModule := entities.NewNeedyVentGasModule(rng)
	ventGasModule.SetBomb(bomb)

	ventGasModule.SetState(entities.NeedyVentGasState{
		DisplayedQuestion:  "vent gas?",
		CountdownStartedAt: time.Now().Unix(),
		CountdownDuration:  30,
	})

	ventGasModuleActor := actors.NewNeedyVentGasModuleActor(ventGasModule)
	ventGasModuleActor.Start()
	defer ventGasModuleActor.Stop()

	sessionID := uuid.New()
	bombID := uuid.New()
	moduleID := ventGasModule.GetModuleID()

	// Answer correctly multiple times
	for i := 0; i < 3; i++ {
		module := ventGasModuleActor.GetModule().(*entities.NeedyVentGasModule)
		question := module.GetCurrentQuestion()

		// Answer correctly based on question
		correctAnswer := question == "vent gas?"

		cmd := &command.NeedyVentGasCommand{
			BaseModuleInputCommand: command.BaseModuleInputCommand{
				SessionID: sessionID,
				BombID:    bombID,
				ModuleID:  moduleID,
			},
			Input: correctAnswer,
		}

		respChan := make(chan actors.Response, 1)
		ventGasModuleActor.Send(actors.ModuleCommandMessage{
			Command:         cmd,
			ResponseChannel: respChan,
		})

		var resp actors.Response
		select {
		case resp = <-respChan:
		case <-time.After(1 * time.Second):
			t.Fatalf("Timeout waiting for response on iteration %d", i)
		}

		if successResp, ok := resp.(actors.SuccessResponse); ok {
			result, ok := successResp.Data.(*command.NeedyVentGasCommandResult)
			assert.True(t, ok, "Expected NeedyVentGasCommandResult type")
			// Needy modules never mark as solved
			assert.False(t, result.Solved, "Needy module should never be marked as solved")
		}
	}
}
