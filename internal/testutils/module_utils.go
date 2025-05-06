package testutils

import (
	"testing"

	"github.com/ZaneH/keep-talking/internal/actors"
	"github.com/ZaneH/keep-talking/internal/domain/entities"
)

func NewTestSimpleWireModule(t *testing.T) *actors.SimpleWiresModuleActor {
	t.Helper()

	module := entities.NewSimpleWiresModule()

	return actors.NewSimpleWireModuleActor(module)
}

func NewTestPasswordModule(t *testing.T, solution string) *actors.PasswordModuleActor {
	t.Helper()

	module := entities.NewPasswordModule(solution)

	return actors.NewPasswordModuleActor(module)
}
