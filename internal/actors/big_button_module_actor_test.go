package actors_test

import (
	"log"
	"testing"
	"time"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBigButtonModuleActor_BasicSolve(t *testing.T) {
	// Arrange
	bomb := entities.NewBomb(valueobject.NewDefaultBombConfig())
	buttonModule := entities.NewBigButtonModule(bomb)
	buttonModuleActor := actors.NewBigButtonModuleActor(bomb, buttonModule)
	buttonModuleActor.Start() // Start the actor to process messages
	defer buttonModuleActor.Stop()

	var specifiedModule *entities.BigButtonModule
	if module, ok := buttonModuleActor.GetModule().(*entities.BigButtonModule); ok {
		specifiedModule = module
	} else {
		t.Fatal("Could not cast to BigButtonModule")
	}

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
		Action: valueobject.PressTypeTap,
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
		t.Logf("Success response received: %+v", successResp.Error())
		assert.True(t, ok, "Expected SuccessResponse type")

		result, ok := successResp.Data.(*command.BigButtonInputCommandResult)
		assert.True(t, ok, "Expected BigButtonInputCommandResult type")

		assert.Equal(t, specifiedModule.IsSolved(), result.Solved, "Solved state should match module state")
		assert.False(t, result.Strike, "No strike should be issued for successful operations")
	}

	assert.Nil(t, resp.Error(), "Expected no error in response")

	log.Printf("Final state: %v", buttonModuleActor.GetModule())
}
