package services

import (
	"context"
	"testing"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/ZaneH/keep-talking/internal/testutils"
)

func TestGameService_CreateGameSession(t *testing.T) {
	// Arrange
	actorSystem := actors.NewActorSystem()
	gameService := NewGameService(actorSystem)
	cmd := &command.CreateGameCommand{}

	// Act
	session, err := gameService.CreateGameSession(context.Background(), cmd)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if session == nil {
		t.Fatal("expected a game session, got nil")
	}
}

func TestGameService_ProcessModuleInput(t *testing.T) {
	// Arrange
	actorSystem := actors.NewActorSystem()
	gameService := NewGameService(actorSystem)

	createCmd := &command.CreateGameCommand{}
	session, err := gameService.CreateGameSession(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("failed to create game session: %v", err)
	}

	wireModule := testutils.NewTestSimpleWireModule(t)
	sessionActor, err := gameService.GetGameSession(context.Background(), session.GetSessionID())
	if err != nil {
		t.Fatalf("failed to get game session actor")
	}

	sessionActor.AddModule(wireModule, valueobject.ModulePosition{
		Row:    0,
		Column: 2,
		Face:   valueobject.Front,
	})

	cmd := &command.SimpleWireInputCommand{
		BaseModuleInputCommand: command.BaseModuleInputCommand{
			SessionID: session.GetSessionID(),
			ModulePosition: valueobject.ModulePosition{
				Row:    0,
				Column: 2,
				Face:   valueobject.Front,
			},
		},
		WireIndex: 0,
	}

	// Act
	result, err := gameService.ProcessModuleInput(context.Background(), cmd)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatal("expected a result, got nil")
	}
}
