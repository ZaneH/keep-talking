package actors

import (
	"context"
	"errors"
	"sync"

	"github.com/ZaneH/keep-talking/internal/application/command"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type GameSessionActor struct {
	sessionID    string
	session      *entities.GameSession
	moduleActors map[valueobject.ModulePosition]ModuleActor
	mu           sync.RWMutex
}

func NewGameSessionActor(sessionId uuid.UUID) *GameSessionActor {
	return &GameSessionActor{
		sessionID:    sessionId.String(),
		session:      entities.NewGameSession(sessionId.String()),
		moduleActors: make(map[valueobject.ModulePosition]ModuleActor),
	}
}

func (g *GameSessionActor) ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	switch c := cmd.(type) {
	case *command.CutWireCommand:
		moduleActor, exists := g.moduleActors[c.ModulePosition]
		if !exists {
			return nil, errors.New("module not found")
		}
		return moduleActor.ProcessCommand(ctx, c)
	case *command.SubmitPasswordCommand:
		moduleActor, exists := g.moduleActors[c.ModulePosition]
		if !exists {
			return nil, errors.New("module not found")
		}
		return moduleActor.ProcessCommand(ctx, c)
	default:
		return nil, errors.New("unknown command type")
	}
}
