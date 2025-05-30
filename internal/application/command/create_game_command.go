package command

import (
	"github.com/google/uuid"
)

type CreateGameCommand struct {
	Seed string
}

type CreateGameCommandResult struct {
	SessionID uuid.UUID
}
