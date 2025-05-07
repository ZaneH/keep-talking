package ports

import (
	"github.com/ZaneH/keep-talking/internal/domain/entities"
	"github.com/google/uuid"
)

type GameSessionManager interface {
	GetGameSession(sessionID uuid.UUID) (SessionActor, error)
}

type SessionActor interface {
	AddBomb(bomb *entities.Bomb) error
}
