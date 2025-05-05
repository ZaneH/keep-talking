package actors

import (
	"context"

	"github.com/ZaneH/keep-talking/internal/application/common"
	"github.com/google/uuid"
)

type ModuleActor interface {
	// Returns the ID of the underlying module.
	GetModuleID() uuid.UUID
	// Returns the underlying module.
	GetModule() common.Module
	// Process a command and returns the result.
	ProcessCommand(ctx context.Context, cmd interface{}) (interface{}, error)
}
