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
	session           *entities.GameSession
	mu                sync.RWMutex
	modules           map[uuid.UUID]ModuleActor
	modulesByPosition map[valueobject.ModulePosition]uuid.UUID
}

func NewGameSessionActor(sessionID uuid.UUID) *GameSessionActor {
	return &GameSessionActor{
		session:           entities.NewGameSession(sessionID),
		modules:           make(map[uuid.UUID]ModuleActor),
		modulesByPosition: make(map[valueobject.ModulePosition]uuid.UUID),
	}
}

func (g *GameSessionActor) AddModule(module ModuleActor, position valueobject.ModulePosition) {
	g.mu.Lock()
	defer g.mu.Unlock()

	moduleId := module.GetModuleID()
	g.modules[moduleId] = module
	g.modulesByPosition[position] = moduleId
}

func (g *GameSessionActor) GetSessionID() uuid.UUID {
	return g.session.SessionID
}

func (g *GameSessionActor) ProcessCommand(ctx context.Context, cmd command.ModuleInputCommand) (interface{}, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	var moduleId uuid.UUID
	var position = cmd.GetModulePosition()

	moduleId, exists := g.modulesByPosition[position]
	if !exists {
		return nil, errors.New("module not found")
	}

	moduleActor, exists := g.modules[moduleId]
	if !exists {
		return nil, errors.New("module actor not found")
	}

	return moduleActor.ProcessCommand(ctx, cmd)
}
