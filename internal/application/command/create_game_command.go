package command

import (
	"github.com/ZaneH/keep-talking/internal/domain/valueobject"
	"github.com/google/uuid"
)

type CreateGameCommand struct {
	Config valueobject.GameConfig
}

type CreateGameCommandResult struct {
	SessionID uuid.UUID
}
