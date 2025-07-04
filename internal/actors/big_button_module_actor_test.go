package actors_test

import (
	"log"
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

func TestBigButtonModuleActor_TwoBatteriesAndDetonate(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	bomb.Batteries = 2
	buttonModule := entities.NewBigButtonModule(rng)
	buttonModule.SetBomb(bomb)

	testState := entities.NewButtonState(rng)
	testState.ButtonColor = valueobject.Red
	testState.Label = "Detonate"
	buttonModule.SetState(testState)

	buttonModuleActor := actors.NewBigButtonModuleActor(buttonModule)
	buttonModuleActor.Start() // Start the actor to process messages
	defer buttonModuleActor.Stop()

	sessionID := uuid.New()
	bombID := bomb.ID
	moduleID := buttonModule.ModuleID

	// Act
	cmd := &command.BigButtonInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		PressType: valueobject.PressTypeTap,
	}

	respChan := make(chan actors.Response, 1)

	buttonModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	// Assert
	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for response")
	}

	if resp.IsSuccess() {
		successResp, ok := resp.(actors.SuccessResponse)
		t.Logf("Success response received: %+v", successResp.Data)
		assert.True(t, ok, "Expected SuccessResponse type")

		result, ok := successResp.Data.(*command.BigButtonInputCommandResult)
		assert.True(t, ok, "Expected BigButtonInputCommandResult type")

		assert.False(t, result.Strike, "No strike should be issued for successful operations")
	}

	assert.Nil(t, resp.Error(), "Expected no error in response")

	log.Printf("Final state: %v", buttonModuleActor.GetModule())
}

func TestBigButtonModuleActor_FRKLitAndThreeBatteries(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	bomb.Batteries = 3
	bomb.Indicators["FRK"] = valueobject.Indicator{Lit: true}

	buttonModule := entities.NewBigButtonModule(rng)
	buttonModule.SetBomb(bomb)

	testState := entities.NewButtonState(rng)
	testState.ButtonColor = valueobject.White
	testState.Label = "Abort"
	buttonModule.SetState(testState)

	buttonModuleActor := actors.NewBigButtonModuleActor(buttonModule)
	buttonModuleActor.Start() // Start the actor to process messages
	defer buttonModuleActor.Stop()

	sessionID := uuid.New()
	bombID := bomb.ID
	moduleID := buttonModule.ModuleID

	// Act
	cmd := &command.BigButtonInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		PressType: valueobject.PressTypeTap,
	}

	respChan := make(chan actors.Response, 1)

	buttonModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	// Assert
	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for response")
	}

	if resp.IsSuccess() {
		successResp, ok := resp.(actors.SuccessResponse)
		t.Logf("Success response received: %+v", successResp.Data)
		assert.True(t, ok, "Expected SuccessResponse type")

		result, ok := successResp.Data.(*command.BigButtonInputCommandResult)
		assert.True(t, ok, "Expected BigButtonInputCommandResult type")

		assert.False(t, result.Strike, "No strike should be issued for successful operations")
	}

	assert.Nil(t, resp.Error(), "Expected no error in response")

	log.Printf("Final state: %v", buttonModuleActor.GetModule())
}

func TestBigButtonModuleActor_YellowButton(t *testing.T) {
	// Arrange
	rng := services.NewSeededRNGFromString("test")
	bomb := entities.NewBomb(rng, valueobject.NewDefaultBombConfig())
	bomb.Batteries = 1
	bomb.Indicators["FRK"] = valueobject.Indicator{Lit: false}

	buttonModule := entities.NewBigButtonModule(rng)
	buttonModule.SetBomb(bomb)

	testState := entities.NewButtonState(rng)
	testState.ButtonColor = valueobject.Yellow
	testState.Label = "Hold"
	buttonModule.SetState(testState)

	buttonModuleActor := actors.NewBigButtonModuleActor(buttonModule)
	buttonModuleActor.Start() // Start the actor to process messages
	defer buttonModuleActor.Stop()

	sessionID := uuid.New()
	bombID := bomb.ID
	moduleID := buttonModule.ModuleID

	// Act
	cmd := &command.BigButtonInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: sessionID,
			BombID:    bombID,
			ModuleID:  moduleID,
		},
		PressType: valueobject.PressTypeTap,
	}

	respChan := make(chan actors.Response, 1)

	buttonModuleActor.Send(actors.ModuleCommandMessage{
		Command:         cmd,
		ResponseChannel: respChan,
	})

	// Assert
	var resp actors.Response
	select {
	case resp = <-respChan:
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for response")
	}

	if successResp, ok := resp.(actors.SuccessResponse); ok {
		if typedResp, ok := successResp.Data.(*command.BigButtonInputCommandResult); ok {
			assert.True(t, typedResp.Strike, "Expected a strike for this operation")
		} else {
			t.Errorf("Expected command result response, found %T", successResp.Data)
		}
	} else {
		t.Error("Expected success response")
	}

	log.Printf("Final state: %v", buttonModuleActor.GetModule())
}
