package entities

import (
	"time"

	"github.com/google/uuid"
)

type GameSession struct {
	SessionID     uuid.UUID
	GameStartedAt *time.Time
}

func NewGameSession(sessionID uuid.UUID) *GameSession {
	return &GameSession{
		SessionID:     sessionID,
		GameStartedAt: nil,
	}
}
