package services_test

import (
	"context"
	"testing"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/application/services"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

func TestGameService_CreateGameSession(t *testing.T) {
	// Arrange
	actorSystem := actors.NewActorSystem()
	gameService := services.NewGameService(actorSystem)
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
	gameService := services.NewGameService(actorSystem)

	createCmd := &command.CreateGameCommand{}
	session, err := gameService.CreateGameSession(context.Background(), createCmd)
	if err != nil {
		t.Fatalf("failed to create game session: %v", err)
	}

	wireModule := actors.NewSimpleWiresModuleActor()

	sessionActor, err := gameService.GetGameSession(context.Background(), session.GetSessionID())
	if err != nil {
		t.Fatalf("failed to get game session actor")
	}

	sessionActor.AddModule(wireModule, valueobject.ModulePosition{
		Row:    0,
		Column: 2,
		Face:   valueobject.Front,
	})

	cmd := &command.SimpleWiresInputCommand{
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
	res, err := gameService.ProcessModuleInput(context.Background(), cmd)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil {
		t.Fatal("expected a result, got nil")
	}
}
