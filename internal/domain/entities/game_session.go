package entities

import "github.com/google/uuid"

type GameSession struct {
	SessionID uuid.UUID
}

func NewGameSession(sessionID uuid.UUID) *GameSession {
	return &GameSession{
		SessionID: sessionID,
	}
}
