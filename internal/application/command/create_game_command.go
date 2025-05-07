package command

import (
	"github.com/google/uuid"
)

type CreateGameCommand struct {
}

type CreateGameCommandResult struct {
	SessionID uuid.UUID
}
