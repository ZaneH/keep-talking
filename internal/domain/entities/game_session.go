package entities

import (
	"time"

	"github.com/ZaneH/keep-talking/internal/domain/ports"
	"github.com/google/uuid"
)

type GameSession struct {
	SessionID     uuid.UUID
	GameStartedAt *time.Time
	RandomService ports.RandomGenerator
}

func NewGameSession(sessionID uuid.UUID) *GameSession {
	return &GameSession{
		SessionID:     sessionID,
		GameStartedAt: nil,
	}
}

func (g *GameSession) SetRandomGenerator(rng ports.RandomGenerator) {
	g.RandomService = rng
}
