package entities

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
)

type GameSession struct {
	SessionId string
	Modules   map[valueobject.ModulePosition]Module
}

func NewGameSession(sessionId string) *GameSession {
	return &GameSession{
		SessionId: sessionId,
		Modules:   make(map[valueobject.ModulePosition]Module),
	}
}
