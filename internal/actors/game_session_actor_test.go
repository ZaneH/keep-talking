package actors_test

import (
	"testing"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/ZaneH/keep-talking/internal/mocks"
	"github.com/google/uuid"
)

func TestGameSessionActor_CreateGame(t *testing.T) {
	sessionID := uuid.New()
	config := valueobject.GameConfig{}
	sessionActor := actors.NewGameSessionActor(sessionID, config)

	if sessionActor.GetSessionID() != sessionID {
		t.Errorf("expected session ID %v, got %v", sessionID, sessionActor.GetSessionID())
	}

	if len(sessionActor.GetModules()) != 0 {
		t.Errorf("expected no modules, got %d", len(sessionActor.GetModules()))
	}
}

func TestGameSessionActor_AddModule(t *testing.T) {
	sessionID := uuid.New()
	config := valueobject.GameConfig{}
	sessionActor := actors.NewGameSessionActor(sessionID, config)

	mockModule := &mocks.MockModuleActor{}
	position := valueobject.ModulePosition{
		Row:    0,
		Column: 0,
		Face:   valueobject.Front,
	}

	sessionActor.AddModule(mockModule, position)

	if len(sessionActor.GetModules()) != 1 {
		t.Errorf("expected 1 module, got %d", len(sessionActor.GetModules()))
	}

	moduleByPosition, err := sessionActor.GetModuleByPosition(position)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if moduleByPosition.GetModuleID() != mockModule.GetModuleID() {
		t.Errorf("expected module ID %v, got %v", mockModule.GetModuleID(), moduleByPosition.GetModuleID())
	}
}
