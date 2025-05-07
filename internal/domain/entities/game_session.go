package entities

import (
	"time"

	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type GameSession struct {
	SessionID     uuid.UUID
	GameStartedAt *time.Time
	Config        valueobject.GameConfig
}

func NewGameSession(sessionID uuid.UUID, config valueobject.GameConfig) *GameSession {
	return &GameSession{
		SessionID:     sessionID,
		GameStartedAt: nil,
		Config:        config,
	}
}
